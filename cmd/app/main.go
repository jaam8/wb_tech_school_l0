package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jaam8/wb_tech_school_l0/internal/config"
	"github.com/jaam8/wb_tech_school_l0/internal/delivery/http/handlers"
	"github.com/jaam8/wb_tech_school_l0/internal/delivery/http/middlewares"
	"github.com/jaam8/wb_tech_school_l0/internal/ports/adapters/broker"
	"github.com/jaam8/wb_tech_school_l0/internal/ports/adapters/cache"
	"github.com/jaam8/wb_tech_school_l0/internal/ports/adapters/storage"
	"github.com/jaam8/wb_tech_school_l0/internal/service"
	"github.com/jaam8/wb_tech_school_l0/pkg/kafka"
	"github.com/jaam8/wb_tech_school_l0/pkg/logger"
	lrucache "github.com/jaam8/wb_tech_school_l0/pkg/lru-cache"
	"github.com/jaam8/wb_tech_school_l0/pkg/postgres"
	"go.uber.org/zap"
)

// @title Order service API
// @version 1.0
// @description Simple API for wb techschool
// @host localhost:8080
// @BasePath /
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx = context.WithValue(ctx, logger.KeyForLogLevel, cfg.LogLevel)
	ctx, _ = logger.New(ctx)

	cacheCfg := cfg.Cache
	postgresCfg := cfg.Postgres
	appCfg := cfg.Service

	pgClient, err := postgres.New(ctx, postgresCfg)
	if err != nil {
		log.Fatalf("failed to create postgres client: %v", err)
	}

	err = postgres.Migrate(ctx, postgresCfg, cfg.MigrationsPath)
	if err != nil {
		log.Fatalf("failed to migrate postgres: %v", err)
	}

	consumer := kafka.NewReader(ctx, cfg.Kafka, appCfg.KafkaTopic, appCfg.KafkaGroupID)

	inMemoryCache := lrucache.New(
		cacheCfg.Capacity,
		time.Duration(cacheCfg.TTL)*time.Minute,
	)

	postgresAdapter := storage.NewPostgresAdapter(pgClient)
	kafkaAdapter := broker.NewKafkaConsumerAdapter(consumer)
	inMemoryCacheAdapter := cache.NewInMemoryCacheAdapter(inMemoryCache)
	srvc := service.New(inMemoryCacheAdapter, kafkaAdapter, postgresAdapter)
	handler := handlers.NewHandler(srvc)
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		AllowHeaders: "Content-Type",
	}), middlewares.LogMiddleware())

	app.Get("/ping", handlers.Ping)
	apiV1 := app.Group("/api/v1")
	apiV1.Get("/orders/:id", handler.GetOrderByID)

	go func() {
		if err = app.Listen(fmt.Sprintf(":%d", appCfg.Port)); err != nil {
			log.Fatalf("failed to start app: %v", err)
		}
	}()

	go srvc.HandleOrdersEvents(ctx, appCfg.BatchSize,
		time.Second*time.Duration(appCfg.FlushTimeout),
	)
	inMemoryCache.StartCleanup(ctx, time.Duration(cacheCfg.CleanupInterval)*time.Minute)

	<-ctx.Done()
	err = app.Shutdown()
	if err != nil {
		logger.Fatal(ctx, "failed to shutdown server", zap.Error(err))
	}
	logger.Info(ctx, "server stopped")
}
