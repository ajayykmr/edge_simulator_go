package kafka

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

var producer sarama.SyncProducer

func InitKafkaProducer(brokers []string) error {
	var err error

	producer, err = connectProducer(brokers)
	if err != nil {
		return fmt.Errorf("failed to initialize Kafka producer: %w", err)
	}

	return nil
}

func CloseKafkaProducer() {
	if producer != nil {
		err := producer.Close()
		if err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
		} else {
			log.Println("Kafka producer closed")
		}
	}
}

func connectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}

func PushDataToKafka(topic string, message []byte) (partition int32, offset int64, err error) {

	if producer == nil {
		return 0, 0, fmt.Errorf("kafka producer is not initialized")
	}

	//Create a new message to send to the Kafka topic
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
		// Value: sarama.StringEncoder(message),
	}

	//Send the message to the Kafka topic
	partition, offset, err = producer.SendMessage(msg)
	if err != nil {
		return 0, 0, err
	}

	return partition, offset, nil
}
