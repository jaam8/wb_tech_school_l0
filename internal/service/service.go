package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jaam8/wb_tech_school_l0/internal/models"
	"github.com/jaam8/wb_tech_school_l0/internal/ports"
	errs "github.com/jaam8/wb_tech_school_l0/pkg/errors"
	"github.com/jaam8/wb_tech_school_l0/pkg/logger"
	"go.uber.org/zap"
)

type Service struct {
	cache   ports.CacheAdapter
	broker  ports.BrokerAdapter
	storage ports.StorageAdapter
}

func New(
	cache ports.CacheAdapter,
	broker ports.BrokerAdapter,
	storage ports.StorageAdapter,
) *Service {
	return &Service{
		cache:   cache,
		broker:  broker,
		storage: storage,
	}
}

func (s *Service) HandleOrdersEvents(ctx context.Context, batchSize int, flushTimeout time.Duration) {
	ticker := time.NewTicker(flushTimeout)
	defer ticker.Stop()

	var batch []*models.Order

	flushBatch := func() {
		if len(batch) == 0 {
			return
		}
		if err := s.storage.SaveOrders(ctx, batch...); err != nil {
			logger.Error(ctx, "failed to save orders batch to storage",
				zap.Error(err),
			)
		}
		logger.Info(ctx, "saved orders batch to storage", zap.Int("count", len(batch)))
		batch = nil
	}

	for {
		select {
		case <-ctx.Done():
			flushBatch()
			logger.Info(ctx, "stop handling kafka consumer")
			return
		case <-ticker.C:
			flushBatch()
		default:
			event, err := s.broker.ConsumeOrderEvent(ctx)
			if err != nil {
				logger.Error(ctx, "failed to consume order event",
					zap.Error(err),
				)
				continue
			}
			if event == nil {
				logger.Error(ctx, "empty order event")
				continue
			}
			if err = event.Validate(); err != nil {
				logger.Warn(ctx, "failed to validate order event",
					zap.String("order_uid", event.OrderUID),
					zap.Error(err),
				)
				continue
			}

			batch = append(batch, event)

			if len(batch) >= batchSize {
				flushBatch()
			}
		}
	}
}

func (s *Service) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	logger.With(ctx,
		zap.String("order_uid", id),
	)

	if id == "" {
		return nil, errs.ErrEmptyOrderUID
	}
	logger.Info(ctx, "get order")

	order, err := s.cache.GetOrder(id)
	if err != nil {
		logger.Warn(ctx, "failed to get order from cache", zap.Error(err))

		order, err = s.storage.GetOrder(ctx, id)
		if err != nil {
			logger.Error(ctx, "failed to get order from storage", zap.Error(err))
			return nil, fmt.Errorf("failed to get order: %w", err)
		}
		err = s.cache.SaveOrder(id, order)
		if err != nil {
			logger.Error(ctx, "failed to save order to cache", zap.Error(err))
		}
	}

	logger.Info(ctx, "got order")
	return order, nil
}
