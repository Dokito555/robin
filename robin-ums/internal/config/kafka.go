package config

import (
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewKafkaProducer(config *viper.Viper, log *logrus.Logger) (sarama.SyncProducer, error) {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Producer.Timeout = 10 * time.Second
	conf.Producer.Retry.Max = 5
	conf.Producer.RequiredAcks = sarama.WaitForAll

	brokers := strings.Split(config.GetString("KAFKA_BROKERS"), ",")
	log.Infof("Kafka brokers: %v", brokers)

	producer, err := sarama.NewSyncProducer(brokers, conf)
	if err != nil {
		log.Errorf("failed to comunicate with kafka brokers: %v", err)
		return nil, err
	}

	return producer, nil
}

func ProduceKafkaMessage(producer sarama.SyncProducer, log *logrus.Logger, topic string, data []byte) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Errorf("Failed to produce message to Kafka: %v", err)
		return err
	}

	log.Infof("Successfully produced message on topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}