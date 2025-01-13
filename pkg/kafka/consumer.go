package kafka

import (
	"context"
	"log"
	"time"
)

func Consume(ctx context.Context, topics []string, maxRetries int, timeout time.Duration) {
	// Subscribe to the provided topics
	err := consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %v", err)
		return
	}

	for {
		log.Println("Kafka consumer running")
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopped")
			return
		default:
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Error consuming message: %v\n", err)
				continue
			}

			// Process the message
			log.Printf("Hello Consumed message from topic %s: %s\n", *msg.TopicPartition.Topic, msg.Value)
		}
	}
}
