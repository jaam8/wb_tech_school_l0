package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jaam8/wb_tech_school_l0/internal/config"
	"github.com/jaam8/wb_tech_school_l0/internal/models"
	"github.com/jaam8/wb_tech_school_l0/internal/ports/adapters/broker"
	"github.com/jaam8/wb_tech_school_l0/pkg/kafka"
	"github.com/jaam8/wb_tech_school_l0/pkg/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx = context.WithValue(ctx, logger.KeyForLogLevel, cfg.LogLevel)
	ctx, _ = logger.New(ctx)

	kafkaCfg := cfg.Kafka
	appCfg := cfg.Service

	kafkaProducer := kafka.NewWriter(ctx, kafkaCfg, appCfg.KafkaTopic)
	defer kafkaProducer.Close()

	err = kafka.CreateTopicWithRetry(
		cfg.Kafka,
		appCfg.KafkaTopic,
		appCfg.KafkaNumPartitions,
		appCfg.KafkaReplicationFactor,
		appCfg.MaxRetries,
	)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create topic: %w", err))
	}

	kafkaAdapter := broker.NewKafkaProducerAdapter(kafkaProducer)

	orders := make([]models.Order, 0, cfg.FakeOrdersCount)

	for range cfg.FakeOrdersCount {
		order := models.GenerateFakeOrder()
		orders = append(orders, order)
	}

	err = kafkaAdapter.SendOrder(ctx, orders...)
	if err != nil {
		log.Println(err)
	}

	log.Println("generate 100 orders, server stopped")
}
