package services

import (
	"sync"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type MessagingService struct {
	Log      *logrus.Logger
	Consumer sarama.Consumer
}

func NewMessagingService(log *logrus.Logger, consumer sarama.Consumer) *MessagingService {
	return &MessagingService{
		Log:      log,
		Consumer: consumer,
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
