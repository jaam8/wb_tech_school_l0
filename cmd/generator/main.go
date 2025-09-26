package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jaam8/wb_tech_school_l0/internal/config"
	"github.com/jaam8/wb_tech_school_l0/internal/models"
	"github.com/jaam8/wb_tech_school_l0/internal/ports/adapters/broker"
	"github.com/jaam8/wb_tech_school_l0/pkg/kafka"
	"github.com/jaam8/wb_tech_school_l0/pkg/logger"
)

func main() {
	ctx := context.Background()
	//ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	//defer stop()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx = context.WithValue(ctx, "log_level", cfg.LogLevel)
	ctx, _ = logger.New(ctx)

	kafkaCfg := cfg.Kafka
	appCfg := cfg.Service

	fmt.Println(kafkaCfg)

	kafkaProducer := kafka.NewWriter(ctx, kafkaCfg, appCfg.KafkaTopic)
	defer kafkaProducer.Close()

	err = kafka.CreateTopicWithRetry(
		cfg.Kafka,
		appCfg.KafkaTopic,
		1,
		1)
	//appCfg.KafkaNumPartitions,
	//appCfg.KafkaReplicationFactor)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create topic: %w", err))
	}

	kafkaAdapter := broker.NewKafkaProducerAdapter(kafkaProducer)

	orders := make([]models.Order, 0, 100)

	for i := 0; i < 100; i++ {
		delivery := models.Delivery{
			Name:    gofakeit.Name(),
			Phone:   "+" + gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Street(),
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		}

		start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 9, 30, 23, 59, 59, 0, time.UTC)
		orderedAt := gofakeit.DateRange(start, end)

		payment := models.Payment{
			Transaction: strings.ReplaceAll(gofakeit.UUID(), "-", ""),
			RequestId:   ".",
			Provider:    gofakeit.Company(),
			Currency:    gofakeit.CurrencyShort(),
			PaymentDt:   int(orderedAt.Add(time.Minute).Unix()),
			Bank: gofakeit.RandomString([]string{
				"sber", "alpha", "vtb", "tinkoff",
				"yapay", "wbpay", "ozon"}),
			DeliveryCost: gofakeit.Number(0, 500),
			GoodsTotal:   gofakeit.Number(1, 10000),
			CustomFee:    gofakeit.Number(0, 10),
		}
		payment.Amount = payment.DeliveryCost + payment.GoodsTotal

		trackNumber := gofakeit.LetterN(12)
		items := make([]models.Item, 0)

		for j := 0; j < gofakeit.Number(1, 3); j++ {
			item := models.Item{
				ChrtId:      gofakeit.Number(1, 999999),
				TrackNumber: trackNumber,
				Price:       gofakeit.Number(100, 10000),
				Rid:         strings.ReplaceAll(gofakeit.UUID(), "-", ""),
				Name:        gofakeit.ProductName(),
				Sale:        gofakeit.Number(0, 99),
				Size:        gofakeit.Numerify("#"),
				NmId:        gofakeit.Number(1, 9999999),
				Brand:       gofakeit.Company(),
				Status:      gofakeit.Number(1, 5),
			}
			item.TotalPrice = int(float64(item.Price) * (1 - float64(item.Sale)/100.0))

			items = append(items, item)
		}

		order := models.Order{
			OrderUid:          strings.ReplaceAll(gofakeit.UUID(), "-", ""),
			TrackNumber:       trackNumber,
			Entry:             "WBIL",
			Delivery:          delivery,
			Payment:           payment,
			Items:             items,
			Locale:            strings.ToUpper(gofakeit.CountryAbr()),
			InternalSignature: ".",
			CustomerId:        strings.ReplaceAll(gofakeit.UUID(), "-", ""),
			DeliveryService:   gofakeit.Company(),
			Shardkey:          gofakeit.Numerify("#"),
			SmId:              gofakeit.Number(1, 100),
			DateCreated:       orderedAt,
			OofShard:          gofakeit.Numerify("#"),
		}

		orders = append(orders, order)
	}

	err = kafkaAdapter.SendOrder(ctx, orders...)
	if err != nil {
		log.Println(err)
	}

	log.Println("generate 100 orders, server stopped")
}
