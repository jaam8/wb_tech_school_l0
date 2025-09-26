package broker

import (
	"context"
	"encoding/json"

	"github.com/jaam8/wb_tech_school_l0/internal/models"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumerAdapter struct {
	consumer *kafka.Reader
}

func NewKafkaConsumerAdapter(consumer *kafka.Reader) *KafkaConsumerAdapter {
	return &KafkaConsumerAdapter{consumer: consumer}
}

func (a KafkaConsumerAdapter) ConsumeOrderEvent(ctx context.Context) (*models.Order, error) {
	var order models.Order
	msg, err := a.consumer.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(msg.Value, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
