package ports

import (
	"context"

	"github.com/jaam8/wb_tech_school_l0/internal/models"
)

type StorageAdapter interface {
	GetOrder(ctx context.Context, id string) (*models.Order, error)
	SaveOrders(ctx context.Context, order ...*models.Order) error
}

type BrokerAdapter interface {
	ConsumeOrderEvent(ctx context.Context) (*models.Order, error)
}

type CacheAdapter interface {
	GetOrder(id string) (*models.Order, error)
	SaveOrder(key string, val *models.Order) error
	//SaveOrders(ctx context.Context, orders ...*models.Order) error
}
