package broker

import (
	"context"
	"encoding/json"

	"github.com/jaam8/wb_tech_school_l0/internal/models"
	"github.com/segmentio/kafka-go"
)

type KafkaProducerAdapter struct {
	producer *kafka.Writer
}

func NewKafkaProducerAdapter(producer *kafka.Writer) *KafkaProducerAdapter {
	return &KafkaProducerAdapter{
		producer: producer,
	}
}

func (a *KafkaProducerAdapter) SendOrder(ctx context.Context, orders ...models.Order) error {
	msgs := make([]kafka.Message, 0, cap(orders))

	for _, o := range orders {
		orderJSON, err := json.Marshal(o)
		if err != nil {
			return err
		}

		msg := kafka.Message{
			Key:   []byte(o.OrderUid),
			Value: orderJSON,
		}

		msgs = append(msgs, msg)
	}

	err := a.producer.WriteMessages(ctx, msgs...)
	return err
}
