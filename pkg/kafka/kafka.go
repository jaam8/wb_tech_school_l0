package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jaam8/wb_tech_school_l0/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Config struct {
	Host    string   `env:"HOST"    env-default:"kafka" yaml:"host"`
	Port    uint16   `env:"PORT"    env-default:"9092"  yaml:"port"`
	Brokers []string `env:"BROKERS" env-separator:","   yaml:"brokers"`

	MinBytes       int `env:"MIN_BYTES"          env-default:"10"      yaml:"min_bytes"`
	MaxBytes       int `env:"MAX_BYTES"          env-default:"1048576" yaml:"max_bytes"` // 1MB
	MaxWaitMs      int `env:"MAX_WAIT_MS"        env-default:"500"     yaml:"max_wait_ms"`
	CommitInterval int `env:"COMMIT_INTERVAL_MS" env-default:"1000"    yaml:"commit_interval_ms"`
}

func NewReader(ctx context.Context, cfg Config, topic, groupID string) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       cfg.MinBytes,
		MaxBytes:       cfg.MaxBytes,
		MaxWait:        time.Duration(cfg.MaxWaitMs) * time.Millisecond,
		CommitInterval: time.Duration(cfg.CommitInterval) * time.Millisecond,
	})
	logger.Info(ctx, "connected to Kafka topic",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("topic", topic),
		zap.String("group_id", groupID),
	)

	return r
}

func NewWriter(ctx context.Context, cfg Config, topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
		Balancer:     &kafka.LeastBytes{},
		Async:        false,
	}

	logger.Info(ctx, "created Kafka writer",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("topic", topic),
	)

	return w
}

func CreateTopicIfNotExists(cfg Config, topic string, numPartitions, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", cfg.Brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp",
		fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return err
	}

	defer controllerConn.Close()

	return controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	})
}

func CreateTopicWithRetry(cfg Config, topic string, numPartitions, replicationFactor, maxRetries int) error {
	var err error
	for i := range maxRetries {
		err = CreateTopicIfNotExists(cfg, topic, numPartitions, replicationFactor)
		if err == nil {
			return nil
		}

		log.Printf("Attempt %d failed: %v\n", i+1, err)
		time.Sleep(time.Second * time.Duration(i))
	}
	return err
}
