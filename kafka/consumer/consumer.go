package kafka

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokers, config)
}

func ConsumerMain() {
	topic := "sensor_data"
	brokers := []string{"localhost:9092"}
	msgCount := 0

	//create a new consumer and start it
	worker, err := ConnectConsumer(brokers)
	if err != nil {
		log.Fatalf("Error connecting to Kafka Consumer: %v", err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Error starting consumer for topic %s: %v", topic, err)
	}

	log.Println("Consumer started, waiting for messages...")

	//Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	//crea a goroutune to run the consumer/worker
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Printf("Error consuming message: %v\n", err)
			case msg := <-consumer.Messages():
				msgCount++
				message := string(msg.Value)
				log.Printf("Received message: %s from topic: %s, partition: %d, offset: %d\n",
					string(message), msg.Topic, msg.Partition, msg.Offset)

			case <-sigChan:
				log.Println("Received shutdown signal, stopping consumer...")
				doneCh <- struct{}{}
			}
		}
	}()

	//Wait for the done channel to be closed
	<-doneCh
	log.Printf("Total messages consumed: %d\n", msgCount)

	if err := consumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	} else {
		log.Println("Consumer closed successfully.")
	}

	if err := worker.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	} else {
		log.Println("Worker closed successfully.")
	}
}
