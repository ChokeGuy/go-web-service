package kafka

import (
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func Consume(topics []string, maxRetries int, timeout time.Duration) *kafka.Message {
	// Subscribe to the provided topics
	err := consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %v", err)
		return nil
	}

	retries := 0
	startTime := time.Now()

	for retries < maxRetries {
		// Check if timeout has been reached
		if time.Since(startTime) > timeout {
			log.Printf("Timeout reached. No message consumed.")
			return nil
		}

		// Read message with a timeout
		msg, err := consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			// Handle timeout errors separately
			if kafkaError, ok := err.(kafka.Error); ok && kafkaError.Code() == kafka.ErrTimedOut {
				retries++
				continue // No message, loop again
			}
			log.Printf("Error consuming message: %v\n", err)
			retries++
			continue
		}

		// Process the message
		log.Printf("Consumed message from topic %s: %s\n", *msg.TopicPartition.Topic, msg.Value)
		return msg
	}

	log.Printf("Max retries reached. No message consumed.")
	return nil
}
