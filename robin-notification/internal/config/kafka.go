package config

import (
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewKafkaConsumer(config *viper.Viper, log *logrus.Logger) (sarama.Consumer, error) {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Producer.Timeout = 10 * time.Second
	conf.Producer.Retry.Max = 5
	conf.Producer.RequiredAcks = sarama.WaitForAll

	brokers := strings.Split(config.GetString("KAFKA_BROKERS"), ",")
	log.Infof("Kafka brokers: %v", brokers)

	consumer, err := sarama.NewConsumer(brokers, conf)
	if err != nil {
		log.Errorf("failed to create Kafka consumer: %v", err)
		return nil, err
	}

	return consumer, nil
}