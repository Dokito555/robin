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