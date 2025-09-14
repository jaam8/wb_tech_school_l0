package ports

import (
	"context"

	"github.com/jaam8/wb_tech_school_l0/internal/models"
)

type StorageAdapter interface {
	GetOrder(ctx context.Context, id string) (*models.Order, error)
	SaveOrder(ctx context.Context, order *models.Order) error
}
