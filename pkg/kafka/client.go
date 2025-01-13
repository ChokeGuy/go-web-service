package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var producer *kafka.Producer
var consumer *kafka.Consumer

func InitProducer(brokers string) error {
	var err error

	// Initialize Producer
	producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	return nil
}

func InitConsumer(brokers string, groupID string, topics []string) error {
	// Initialize Consumer
	var err error
	consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     brokers,
		"group.id":              groupID,
		"auto.offset.reset":     "earliest",
		"enable.auto.commit":    true,
		"session.timeout.ms":    10000,
		"heartbeat.interval.ms": 3000,
		"max.poll.interval.ms":  300000,
		"fetch.min.bytes":       1,
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}

	return nil
}
