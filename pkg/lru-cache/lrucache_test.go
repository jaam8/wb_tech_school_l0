package lrucache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryCache_Set(t *testing.T) {
	tests := []struct {
		name       string
		capacity   int
		ttl        time.Duration
		operations []struct {
			key   interface{}
			value interface{}
		}
		wantErr     bool
		wantLen     int
		description string
	}{
		{
			name:     "set single item",
			capacity: 10,
			ttl:      time.Minute,
			operations: []struct {
				key   interface{}
				value interface{}
			}{
				{key: "key1", value: "value1"},
			},
			wantErr:     false,
			wantLen:     1,
			description: "should successfully set a single item",
		},
		{
			name:     "set multiple items",
			capacity: 10,
			ttl:      time.Minute,
			operations: []struct {
				key   interface{}
				value interface{}
			}{
				{key: "key1", value: "value1"},
				{key: "key2", value: "value2"},
				{key: "key3", value: "value3"},
			},
			wantErr:     false,
			wantLen:     3,
			description: "should successfully set multiple items",
		},
		{
			name:     "update existing item",
			capacity: 10,
			ttl:      time.Minute,
			operations: []struct {
				key   interface{}
				value interface{}
			}{
				{key: "key1", value: "value1"},
				{key: "key1", value: "value2"},
			},
			wantErr:     false,
			wantLen:     1,
			description: "should update value for existing key",
		},
		{
			name:     "evict LRU item when capacity exceeded",
			capacity: 2,
			ttl:      time.Minute,
			operations: []struct {
				key   interface{}
				value interface{}
			}{
				{key: "key1", value: "value1"},
				{key: "key2", value: "value2"},
				{key: "key3", value: "value3"},
			},
			wantErr:     false,
			wantLen:     2,
			description: "should evict least recently used item when capacity is exceeded",
		},
		{
			name:     "set nil value",
			capacity: 10,
			ttl:      time.Minute,
			operations: []struct {
				key   interface{}
				value interface{}
			}{
				{key: "key1", value: nil},
			},
			wantErr:     true,
			wantLen:     0,
			description: "should return error when setting nil value",
		},
		{
			name:     "set with zero capacity",
			capacity: 0,
			ttl:      time.Minute,
			operations: []struct {
				key   interface{}
				value interface{}
			}{
				{key: "key1", value: "value1"},
				{key: "key2", value: "value2"},
			},
			wantErr:     false,
			wantLen:     2,
			description: "should allow unlimited items with zero capacity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := New(tt.capacity, tt.ttl)
			var err error

			for _, op := range tt.operations {
				err = cache.Set(op.key, op.value)
				if !tt.wantErr {
					require.NoError(t, err, "Set() returned unexpected error")
				}
			}

			if tt.wantErr {
				require.Error(t, err, "Set() expected an error but got none")
			} else {
				require.NoError(t, err, "Set() returned unexpected error")
			}
			require.Equal(t, tt.wantLen, cache.list.Len(), "cache length after Set() mismatch")
		})
	}
}

func TestInMemoryCache_Get(t *testing.T) {
	tests := []struct {
		name        string
		capacity    int
		ttl         time.Duration
		setup       func(*InMemoryCache)
		key         interface{}
		wantValue   interface{}
		wantErr     error
		description string
	}{
		{
			name:     "get existing item",
			capacity: 10,
			ttl:      time.Minute,
			setup: func(c *InMemoryCache) {
				c.Set("key1", "value1")
			},
			key:         "key1",
			wantValue:   "value1",
			wantErr:     nil,
			description: "should retrieve existing item",
		},
		{
			name:     "get non-existing item",
			capacity: 10,
			ttl:      time.Minute,
			setup: func(c *InMemoryCache) {
				c.Set("key1", "value1")
			},
			key:         "key2",
			wantValue:   nil,
			wantErr:     ErrNotFound,
			description: "should return ErrNotFound for non-existing key",
		},
		{
			name:     "get expired item",
			capacity: 10,
			ttl:      time.Millisecond,
			setup: func(c *InMemoryCache) {
				c.Set("key1", "value1")
				time.Sleep(2 * time.Millisecond)
			},
			key:         "key1",
			wantValue:   nil,
			wantErr:     ErrExpired,
			description: "should return ErrExpired for expired item",
		},
		{
			name:     "get updates LRU order",
			capacity: 2,
			ttl:      time.Minute,
			setup: func(c *InMemoryCache) {
				c.Set("key1", "value1")
				c.Set("key2", "value2")
				c.Get("key1")
				c.Set("key3", "value3")
			},
			key:         "key1",
			wantValue:   "value1",
			wantErr:     nil,
			description: "should update LRU order on get",
		},
		{
			name:     "get different types",
			capacity: 10,
			ttl:      time.Minute,
			setup: func(c *InMemoryCache) {
				c.Set(123, "numeric key")
				c.Set("string", 456)
			},
			key:         123,
			wantValue:   "numeric key",
			wantErr:     nil,
			description: "should work with different key types",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := New(tt.capacity, tt.ttl)
			if tt.setup != nil {
				tt.setup(cache)
			}

			got, err := cache.Get(tt.key)

			require.ErrorIs(t, tt.wantErr, err, "unexpected error from Get()")
			if err == nil {
				require.Equal(t, tt.wantValue, got, "Get() returned unexpected value")
			}
		})
	}
}

func TestInMemoryCache_Delete(t *testing.T) {
	tests := []struct {
		name        string
		capacity    int
		ttl         time.Duration
		setup       func(*InMemoryCache)
		deleteKey   interface{}
		wantErr     error
		wantLen     int
		description string
	}{
		{
			name:     "delete existing item",
			capacity: 10,
			ttl:      time.Minute,
			setup: func(c *InMemoryCache) {
				c.Set("key1", "value1")
				c.Set("key2", "value2")
			},
			deleteKey:   "key1",
			wantErr:     nil,
			wantLen:     1,
			description: "should successfully delete existing item",
		},
		{
			name:     "delete non-existing item",
			capacity: 10,
			ttl:      time.Minute,
			setup: func(c *InMemoryCache) {
				c.Set("key1", "value1")
			},
			deleteKey:   "key2",
			wantErr:     ErrNotFound,
			wantLen:     1,
			description: "should return ErrNotFound when deleting non-existing item",
		},
		{
			name:        "delete from empty cache",
			capacity:    10,
			ttl:         time.Minute,
			setup:       func(c *InMemoryCache) {},
			deleteKey:   "key1",
			wantErr:     ErrNotFound,
			wantLen:     0,
			description: "should return ErrNotFound on empty cache",
		},
		{
			name:     "delete all items",
			capacity: 10,
			ttl:      time.Minute,
			setup: func(c *InMemoryCache) {
				c.Set("key1", "value1")
				c.Delete("key1")
			},
			deleteKey:   "key1",
			wantErr:     ErrNotFound,
			wantLen:     0,
			description: "should handle multiple deletes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := New(tt.capacity, tt.ttl)
			if tt.setup != nil {
				tt.setup(cache)
			}

			err := cache.Delete(tt.deleteKey)

			require.ErrorIs(t, err, tt.wantErr, "unexpected error from Delete()")
			require.Equal(t, tt.wantLen, cache.list.Len(), "cache length after delete mismatch")
		})
	}
}

func TestInMemoryCache_StartCleanup(t *testing.T) {
	tests := []struct {
		name            string
		capacity        int
		ttl             time.Duration
		cleanupInterval time.Duration
		setup           func(*InMemoryCache)
		waitTime        time.Duration
		wantLenAfter    int
		description     string
	}{
		{
			name:            "cleanup removes expired items",
			capacity:        10,
			ttl:             50 * time.Millisecond,
			cleanupInterval: 30 * time.Millisecond,
			setup: func(c *InMemoryCache) {
				c.Set("key1", "value1")
				c.Set("key2", "value2")
			},
			waitTime:     100 * time.Millisecond,
			wantLenAfter: 0,
			description:  "should remove expired items during cleanup",
		},
		{
			name:            "cleanup keeps non-expired items",
			capacity:        10,
			ttl:             time.Second,
			cleanupInterval: 20 * time.Millisecond,
			setup: func(c *InMemoryCache) {
				c.Set("key1", "value1")
				c.Set("key2", "value2")
			},
			waitTime:     50 * time.Millisecond,
			wantLenAfter: 2,
			description:  "should keep non-expired items during cleanup",
		},
		{
			name:            "cleanup stops on context cancel",
			capacity:        10,
			ttl:             time.Minute,
			cleanupInterval: 10 * time.Millisecond,
			setup: func(c *InMemoryCache) {
				c.Set("key1", "value1")
			},
			waitTime:     20 * time.Millisecond,
			wantLenAfter: 1,
			description:  "should stop cleanup on context cancellation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := New(tt.capacity, tt.ttl)
			if tt.setup != nil {
				tt.setup(cache)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			cache.StartCleanup(ctx, tt.cleanupInterval)
			time.Sleep(tt.waitTime)
			cancel()

			cache.mu.RLock()
			gotLen := cache.list.Len()
			cache.mu.RUnlock()

			require.Equal(t, tt.wantLenAfter, gotLen, "cache length after cleanup mismatch")
		})
	}
}

func TestInMemoryCache_New(t *testing.T) {
	tests := []struct {
		name        string
		capacity    int
		ttl         time.Duration
		description string
	}{
		{
			name:        "create cache with valid parameters",
			capacity:    100,
			ttl:         time.Minute,
			description: "should create cache with specified capacity and TTL",
		},
		{
			name:        "create cache with zero capacity",
			capacity:    0,
			ttl:         time.Minute,
			description: "should create cache with unlimited capacity",
		},
		{
			name:        "create cache with zero TTL",
			capacity:    100,
			ttl:         0,
			description: "should create cache with zero TTL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := New(tt.capacity, tt.ttl)

			require.NotNil(t, cache, "New() returned nil")
			require.Equal(t, tt.capacity, cache.cap, "cache capacity mismatch")
			require.Equal(t, cache.TTL, tt.ttl, "cache TTL mismatch")
			require.NotNil(t, cache.items, "cache.items should not be nil")
			require.NotNil(t, cache.list, "cache.list should not be nil")
		})
	}
}

func TestInMemoryCache_Concurrency(t *testing.T) {
	tests := []struct {
		name        string
		capacity    int
		ttl         time.Duration
		goroutines  int
		operations  int
		description string
	}{
		{
			name:        "concurrent set operations",
			capacity:    100,
			ttl:         time.Minute,
			goroutines:  10,
			operations:  100,
			description: "should handle concurrent Set operations safely",
		},
		{
			name:        "concurrent get operations",
			capacity:    100,
			ttl:         time.Minute,
			goroutines:  10,
			operations:  100,
			description: "should handle concurrent Get operations safely",
		},
		{
			name:        "concurrent mixed operations",
			capacity:    100,
			ttl:         time.Minute,
			goroutines:  10,
			operations:  50,
			description: "should handle concurrent mixed operations safely",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := New(tt.capacity, tt.ttl)

			done := make(chan bool)

			for i := 0; i < tt.goroutines; i++ {
				go func(id int) {
					for j := 0; j < tt.operations; j++ {
						key := id*tt.operations + j
						cache.Set(key, j)
						cache.Get(key)
						if j%2 == 0 {
							cache.Delete(key)
						}
					}
					done <- true
				}(i)
			}

			for i := 0; i < tt.goroutines; i++ {
				<-done
			}

			assert.LessOrEqual(t, len(cache.items), tt.capacity)
		})
	}
}
