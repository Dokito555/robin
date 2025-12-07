package services

import (
	// "github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type MessagingService struct {
	Log      *logrus.Logger
	// Producer sarama.SyncProducer
}

func NewMessagingService(log *logrus.Logger) *MessagingService {
	return &MessagingService{
		Log:      log,
		// Producer: producer,
	}
}

func (m *MessagingService) ProduceKafkaMessage(topic string, data []byte) error {
	// message := &sarama.ProducerMessage{
	// 	Topic: topic,
	// 	Value: sarama.ByteEncoder(data),
	// }

	// partition, offset, err := m.Producer.SendMessage(message)
	// if err != nil {
	// 	m.Log.Errorf("Failed to produce message to Kafka: %v", err)
	// 	return err
	// }

	// m.Log.Infof("Successfully produced message on topic %s, partition %d, offset %d", topic, partition, offset)
	// return nil
	return nil
}
