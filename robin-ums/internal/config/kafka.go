package config

import (
	"fmt"
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

	conf.Net.DialTimeout = 10 * time.Second
	conf.Net.ReadTimeout = 10 * time.Second
	conf.Net.WriteTimeout = 10 * time.Second

	brokers := strings.Split(config.GetString("KAFKA_BROKERS"), ",")
	log.Infof("Kafka brokers: %v", brokers)

	maxRetries := 5
	retryDelay := 5 * time.Second
	
	var producer sarama.SyncProducer
	var err error

	for i := 0; i < maxRetries; i++ {
		producer, err = sarama.NewSyncProducer(brokers, conf)
		if err == nil {
			log.Info("Successfully connected to Kafka brokers")
			return producer, nil
		}

		log.Errorf("Attempt %d/%d: Failed to communicate with kafka brokers: %v", i+1, maxRetries, err)
		
		if i < maxRetries-1 {
			log.Infof("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to Kafka after %d attempts: %w", maxRetries, err)
}