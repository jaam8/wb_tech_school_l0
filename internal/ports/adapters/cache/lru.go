package cache

import (
	"errors"
	"fmt"

	"github.com/jaam8/wb_tech_school_l0/internal/models"
	errs "github.com/jaam8/wb_tech_school_l0/pkg/errors"
	lrucache "github.com/jaam8/wb_tech_school_l0/pkg/lru-cache"
)

type InMemoryCacheAdapter struct {
	client *lrucache.InMemoryCache
}

func NewInMemoryCacheAdapter(client *lrucache.InMemoryCache) *InMemoryCacheAdapter {
	return &InMemoryCacheAdapter{
		client: client,
	}
}

func (a *InMemoryCacheAdapter) GetOrder(key string) (*models.Order, error) {
	val, err := a.client.Get(key)
	if err != nil {
		if errors.Is(err, lrucache.ErrNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}

	order, ok := val.(*models.Order)
	if !ok {
		return nil, fmt.Errorf("invalid type for order")
	}
	return order, nil
}

func (a *InMemoryCacheAdapter) SaveOrder(key string, val *models.Order) error {
	return a.client.Set(key, val)
}
