package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
var kafkaMessagesProduced = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kafka_produced_messages_total",
		Help: "Total number of Kafka messages produced",
	},
	[]string{"topic"},
)

func init() {
	prometheus.MustRegister(kafkaMessagesProduced)
}

func main() {
	// Start Prometheus HTTP server
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting Prometheus metrics server on :2112")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	// Kafka producer setup
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	topic := "test_topic"

	// Produce messages
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Message %d", i)
		err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(message),
		}, nil)

		if err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			kafkaMessagesProduced.WithLabelValues(topic).Inc() // Update Prometheus counter
		}
	}
	// Pause
	time.Sleep(1 * time.Minute)

	// Flush all messages
	producer.Flush(5000)
	fmt.Println("All messages sent successfully!")
}
