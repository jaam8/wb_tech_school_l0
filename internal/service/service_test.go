package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jaam8/wb_tech_school_l0/internal/models"
	errs "github.com/jaam8/wb_tech_school_l0/pkg/errors"
	"github.com/jaam8/wb_tech_school_l0/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCacheAdapter struct {
	mock.Mock
}

func (m *MockCacheAdapter) GetOrder(key string) (*models.Order, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockCacheAdapter) SaveOrder(key string, val *models.Order) error {
	args := m.Called(key, val)
	return args.Error(0)
}

type MockStorageAdapter struct {
	mock.Mock
}

func (m *MockStorageAdapter) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockStorageAdapter) SaveOrders(ctx context.Context, order ...*models.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

type MockBrokerAdapter struct {
	mock.Mock
}

func (m *MockBrokerAdapter) ConsumeOrderEvent(ctx context.Context) (*models.Order, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func TestService_GetOrder(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		want      *models.Order
		wantErr   error
		mockSetup func(storage *MockStorageAdapter, cache *MockCacheAdapter)
	}{
		{
			name: "success got order from cache",
			id:   "test_order_uid",
			want: &models.Order{
				OrderUID: "test_order_uid",
			},
			wantErr: nil,
			mockSetup: func(storage *MockStorageAdapter, cache *MockCacheAdapter) {
				cache.On("GetOrder", "test_order_uid").
					Return(&models.Order{
						OrderUID: "test_order_uid",
					}, nil)
			},
		},
		{
			name: "success got order from storage",
			id:   "test_order_uid",
			want: &models.Order{
				OrderUID: "test_order_uid",
			},
			wantErr: nil,
			mockSetup: func(storage *MockStorageAdapter, cache *MockCacheAdapter) {
				cache.On("GetOrder", "test_order_uid").
					Return(nil, errs.ErrOrderNotFound)
				storage.On("GetOrder", mock.Anything, "test_order_uid").
					Return(&models.Order{
						OrderUID: "test_order_uid",
					}, nil)
				cache.On("SaveOrder", "test_order_uid", &models.Order{OrderUID: "test_order_uid"}).
					Return(nil)
			},
		},
		{
			name:    "order not found",
			id:      "test_order_uid",
			want:    nil,
			wantErr: errs.ErrOrderNotFound,
			mockSetup: func(storage *MockStorageAdapter, cache *MockCacheAdapter) {
				cache.On("GetOrder", "test_order_uid").
					Return(nil, errs.ErrOrderNotFound)
				storage.On("GetOrder", mock.Anything, "test_order_uid").
					Return(nil, errs.ErrOrderNotFound)
			},
		},
		{
			name:    "empty order id",
			id:      "",
			want:    nil,
			wantErr: errs.ErrEmptyOrderUID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := new(MockStorageAdapter)
			cache := new(MockCacheAdapter)
			if tt.mockSetup != nil {
				tt.mockSetup(storage, cache)
			}

			service := New(cache, nil, storage)

			ctx := context.Background()
			ctx, _ = logger.New(ctx)
			order, err := service.GetOrder(ctx, tt.id)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, "test_order_uid", order.OrderUID)
			}

			cache.AssertExpectations(t)
			storage.AssertExpectations(t)
		})
	}
}

func TestService_HandleOrdersEvents(t *testing.T) {
	orders := make([]*models.Order, 0, 2)

	for i := 0; len(orders) < 2; i++ {
		order := models.GenerateFakeOrder()
		if err := order.Validate(); err != nil {
			continue
		}
		orders = append(orders, &order)
	}

	tests := []struct {
		name      string
		batchSize int
		flushTime time.Duration
		timeout   time.Duration
		events    []*models.Order
		mockSetup func(storage *MockStorageAdapter, broker *MockBrokerAdapter)
	}{
		{
			name:      "success save orders batch",
			batchSize: 2,
			flushTime: time.Millisecond * 100,
			timeout:   time.Millisecond * 150,
			events:    orders,
			mockSetup: func(storage *MockStorageAdapter, broker *MockBrokerAdapter) {
				broker.On("ConsumeOrderEvent", mock.Anything).
					Return(orders[0], nil).Once()
				broker.On("ConsumeOrderEvent", mock.Anything).
					Return(orders[1], nil).Once()
				broker.On("ConsumeOrderEvent", mock.Anything).
					Run(func(args mock.Arguments) {
						ctx := args.Get(0).(context.Context)
						select {
						case <-ctx.Done():
							return
						case <-time.After(time.Millisecond * 10):
						}
					}).
					Return(nil, fmt.Errorf("no more events"))

				storage.On("SaveOrders", mock.Anything, orders).Return(nil).Once()
			},
		},
		{
			name:      "flush by timeout",
			batchSize: 10,
			flushTime: time.Millisecond * 50,
			timeout:   time.Millisecond * 100,
			events:    orders[:1],
			mockSetup: func(storage *MockStorageAdapter, broker *MockBrokerAdapter) {
				broker.On("ConsumeOrderEvent", mock.Anything).
					Return(orders[0], nil).Once()
				broker.On("ConsumeOrderEvent", mock.Anything).
					Run(func(args mock.Arguments) {
						ctx := args.Get(0).(context.Context)
						select {
						case <-ctx.Done():
							return
						case <-time.After(time.Millisecond * 10):
						}
					}).
					Return(nil, fmt.Errorf("no more events"))

				storage.On("SaveOrders", mock.Anything, []*models.Order{orders[0]}).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := new(MockStorageAdapter)
			broker := new(MockBrokerAdapter)
			if tt.mockSetup != nil {
				tt.mockSetup(storage, broker)
			}

			service := New(nil, broker, storage)

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()
			ctx, _ = logger.New(ctx)

			done := make(chan struct{})
			go func() {
				defer close(done)
				service.HandleOrdersEvents(ctx, tt.batchSize, tt.flushTime)
			}()

			<-done

			broker.AssertExpectations(t)
			storage.AssertExpectations(t)
		})
	}
}
