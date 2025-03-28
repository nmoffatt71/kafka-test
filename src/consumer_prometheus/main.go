package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
var kafkaMessagesConsumed = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kafka_consumed_messages_total",
		Help: "Total number of Kafka messages consumed",
	},
	[]string{"topic"},
)

func init() {
	prometheus.MustRegister(kafkaMessagesConsumed)
}

func main() {
	// Start Prometheus HTTP server
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting Prometheus metrics server on :2113")
		log.Fatal(http.ListenAndServe(":2113", nil))
	}()

	// Kafka consumer setup
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	topic := "test_topic"
	consumer.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received: %s\n", string(msg.Value))
			kafkaMessagesConsumed.WithLabelValues(topic).Inc() // Update Prometheus counter
		} else {
			fmt.Printf("Consumer error: %v\n", err)
		}
	}
}
