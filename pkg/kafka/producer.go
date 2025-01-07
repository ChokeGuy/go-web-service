package kafka

import (
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

func Produce(topic string, messages []string) {
	for _, msg := range messages {
		// Chuẩn bị message Kafka
		kafkaMsg := &kafka.Message{
			Key:            []byte(uuid.NewString()),
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(msg),
			Timestamp:      time.Now(),
		}

		// Gửi message đến Kafka
		err := producer.Produce(kafkaMsg, nil)
		if err != nil {
			log.Printf("Error producing message to topic %s: %v\n", topic, err)
		} else {
			log.Printf("Produced message to topic %s: %s\n", topic, msg)
		}
	}

	// Đợi tất cả message trong buffer được gửi
	producer.Flush(15 * 1000) // Thời gian chờ tối đa là 15 giây
}
