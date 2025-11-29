package services

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Dokito555/robin-notification/internal/model"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type MessagingService struct {
	Log                 *logrus.Logger
	Consumer            sarama.Consumer
	Notifier			INotifier 
}

func NewMessagingService(log *logrus.Logger, consumer sarama.Consumer) *MessagingService {
	return &MessagingService{
		Log:                 log,
		Consumer:            consumer,
	}
}

func (m *MessagingService) ConsumeKafkaMessage(topic string) (chan []byte, chan error, func(), error) {
	partitions, err := m.Consumer.Partitions(topic)
	if err != nil {
		m.Log.Errorf("failed to get partitions for topic %s: %v", topic, err)
		return nil, nil, nil, err
	}

	// create a wait group to handle concurrent partition consumption
	var wg sync.WaitGroup

	// create channels for messages and errors
	messages := make(chan []byte)
	errors := make(chan error)

	// create a channel to signal shutdown
	shutdown := make(chan struct{})

	// start consumer for each partition
	for _, partition := range partitions {
		wg.Add(1)

		go func(partition int32) {
			defer wg.Done()

			// create a partition consumer
			partitionConsumer, err := m.Consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				m.Log.Errorf("failed to create partition consumer for topic %s, partition %d: %v", topic, partition, err)
				return
			}
			defer partitionConsumer.Close()

			m.Log.Infof("started consuming from topic %s, partition %d", topic, partition)

			// consume messages
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					m.Log.Infof("recieved message from topic %s, partition %d, offset %d", msg.Topic, msg.Partition, msg.Offset)

					// process the message with the provided handler
					// if err := handler(msg.Value); err != nil {
					// 	m.Log.Errorf("error processing message: %v", err)
					// }

					// send the message to the channel
					messages <- msg.Value

				case err := <-partitionConsumer.Errors():
					m.Log.Errorf("error consuming from topic %s, partition %d: %v", topic, partition, err)

				case <-shutdown:
					m.Log.Infof("Shutting down consumer for topic %s, partition %d", topic, partition)
					return
				}
			}
		}(partition)
	}

	// return a closure function that can be called to gracefully shut down the consumer
	closeFunc := func() {
		close(shutdown)
		wg.Wait()
		close(messages)
		close(errors)
	}

	// return a function that can be called to gracefully shut down the consumer
	return messages, errors, closeFunc, nil
}

func GracefulShutdownConsumer(consumer sarama.Consumer, shutdownChan chan struct{}, wg *sync.WaitGroup, log *logrus.Logger) {
	// close the shutdown channel to signal all partition consumers to stop
	close(shutdownChan)

	// wait for all partition consumers to finish
	wg.Wait()

	// close the consumer
	if err := consumer.Close(); err != nil {
		log.Errorf("error closing consumer: %v", err)
	}

	log.Info("kafka consumer has been gracefully shut down")
}

func (m *MessagingService) SetNotifier(n INotifier) {
	m.Notifier = n
}

// HandleKafkaMessage is a default handler function that unmarshals the Kafka message and triggers email sending.
// It assumes that the message JSON contains "username", "role", and "recipient" fields.
func (m *MessagingService) HandleKafkaMessage(msg []byte) error {
	var registerMessage struct {
		Username  string `json:"username"`
		Role      string `json:"role"`
		Recipient string `json:"recipient"`
	}

	if err := json.Unmarshal(msg, &registerMessage); err != nil {
		m.Log.Errorf("Failed to unmarshal Kafka message: %v", err)
		return err
	}

	// Build the internal notification request.
	req := &model.InternalNotificationRequest{
		Recipient: registerMessage.Recipient,
		Placeholder: map[string]interface{}{
			"username": registerMessage.Username,
			"role":     registerMessage.Role,
		},
	}

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Immediately trigger the email sending process.
	if err := m.Notifier.SendEmail(ctx, req); err != nil {
		m.Log.Errorf("Failed to send email: %v", err)
		return err
	}

	m.Log.Infof("Email triggered successfully for recipient %s", registerMessage.Recipient)
	return nil
}
