package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Create a Kafka producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	topic := "quickstart-events"

	// Produce a message
	message := kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("Hello Kafka from GO. Updated!"),
	}

	err = producer.Produce(&message, nil)
	if err != nil {
		log.Printf("Failed to produce message: %s", err)
	}

	// Wait for all messages to be delivered
	producer.Flush(15000)
	producer.Close()

	fmt.Println("Message sent successfully")
}
