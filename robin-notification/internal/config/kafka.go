package config

import (
	"strings"
	"sync"
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
	log.Infof("kafka brokers: %v", brokers)

	consumer, err := sarama.NewConsumer(brokers, conf)
	if err != nil {
		log.Errorf("failed to create Kafka consumer: %v", err)
		return nil, err
	}

	return consumer, nil
}

func ConsumeKafkaMessage(consumer sarama.Consumer, log *logrus.Logger, topic string, handler func([]byte) error) error {
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Errorf("failed to get partitions for topic %s: %v", topic, err)
		return err
	}

	// create a wait group to handle concurrent partition consumption
	var wg sync.WaitGroup

	// create a channel to signal shutdown
	shutdown := make(chan struct{})

	// start consumer for each partition
	for _, partition := range partitions {
		wg.Add(1)

		go func(partition int32) {
			defer wg.Done()

			// create a partition consumer
			partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Errorf("failed to create partition consumer for topic %s, partition %d: %v", topic, partition, err)
				return
			}
			defer partitionConsumer.Close()

			log.Infof("started consuming from topic %s, partition %d", topic, partition)

			// consume messages
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					log.Infof("recieved message from topic %s, partition %d, offset %d", msg.Topic, msg.Partition, msg.Offset)

					// process the message with the provided handler
					if err := handler(msg.Value); err != nil {
						log.Errorf("error processing message: %v", err)
					}

				case err := <-partitionConsumer.Errors():
					log.Errorf("error consuming from topic %s, partition %d: %v", topic, partition, err)

				case <- shutdown:
					log.Infof("Shutting down consumer for topic %s, partition %d", topic, partition)
					return
				}
			}
		}(partition)
	}
	// return a function that can be called to gracefully shut down the consumer
	return nil
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