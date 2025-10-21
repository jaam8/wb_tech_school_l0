package lrucache

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"
)

// InMemoryCache implements a thread-safe LRU cache with TTL support
type InMemoryCache struct {
	mu    sync.RWMutex
	items map[interface{}]*list.Element
	list  *list.List
	cap   int
	TTL   time.Duration
}

// item represents a single cache entry
type item struct {
	key       interface{}
	value     interface{}
	expiredAt time.Time
}

// New creates and returns a new InMemoryCache instance
//   - cap: max number of items the cache can hold
//   - ttl: time to live for duration for each cached item
func New(cap int, ttl time.Duration) *InMemoryCache {
	return &InMemoryCache{
		items: make(map[interface{}]*list.Element),
		list:  list.New(),
		cap:   cap,
		TTL:   ttl,
	}
}

// Set inserts or updates a key-value pair in the cache
//
// If the key already exists, its value and their TTL are updated.
// When the cache exceeds its capacity, the least recently used item is removed
func (c *InMemoryCache) Set(key, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		itm, ok := elem.Value.(*item)
		if !ok {
			return ErrUnexpectedType
		}
		itm.value = value
		itm.expiredAt = time.Now().Add(c.TTL)
		return nil
	}

	if value == nil {
		return fmt.Errorf("cannot set nil value to key %v", key)
	}

	itm := &item{
		key:       key,
		value:     value,
		expiredAt: time.Now().Add(c.TTL),
	}
	elem := c.list.PushFront(itm)
	c.items[key] = elem

	if c.cap > 0 && c.list.Len() > c.cap {
		last := c.list.Back()
		if last != nil {
			c.list.Remove(last)
			lastItem, ok := last.Value.(*item)
			if !ok {
				return ErrUnexpectedType
			}
			delete(c.items, lastItem.key)
		}
	}

	return nil
}

// Get returns the value associated with the given key.
// If the key is expired, returns ErrExpired and ErrNotFound if the key is not found
func (c *InMemoryCache) Get(key interface{}) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.items[key]
	if !ok {
		return nil, ErrNotFound
	}

	itm, ok := elem.Value.(*item)
	if !ok {
		return nil, ErrUnexpectedType
	}

	if itm.value == nil {
		return nil, ErrNotFound
	}

	if time.Now().After(itm.expiredAt) {
		return nil, ErrExpired
	}

	c.list.MoveToFront(elem)

	return itm.value, nil
}

// Delete a key from the cache.
// If the key is not found, return ErrNotFound
func (c *InMemoryCache) Delete(key interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.Remove(elem)
		delete(c.items, key)
		return nil
	}

	return ErrNotFound
}

// StartCleanup launches a background goroutine that periodically
// removes expired items from the cache.
//   - interval: cleaning interval for old items
func (c *InMemoryCache) StartCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.mu.Lock()
				for key, elem := range c.items {
					itm, ok := elem.Value.(*item)
					if !ok {
						continue
					}
					if time.Now().After(itm.expiredAt) {
						c.list.Remove(elem)
						delete(c.items, key)
					}
				}
				c.mu.Unlock()

			case <-ctx.Done():
				return
			}
		}
	}()
}
