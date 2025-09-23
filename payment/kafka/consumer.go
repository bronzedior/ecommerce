package kafka

import (
	"context"
	"encoding/json"
	"log"
	"payment/models"

	"github.com/segmentio/kafka-go"
)

func StartOrderConsumer(broker string, topic string, handler func(models.OrderCreatedEvent)) {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "paymentfc",
	})

	go func(r *kafka.Reader) {
		for {
			message, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Println("Error Read Message Kafka: ", err.Error())
				continue
			}

			var event models.OrderCreatedEvent
			err = json.Unmarshal(message.Value, &event)
			if err != nil {
				log.Println("Error Unmarshal Message: ", err.Error())
				continue
			}

			log.Printf("Received Event Order Created: %+v", event)
			handler(event)
		}
	}(consumer)
}
