package order

import (
	"context"
	"fmt"
	"sync"
	"time"

	appPorts "github.com/D1sordxr/wb-tech-l0/internal/domain/app/ports"
	"github.com/D1sordxr/wb-tech-l0/internal/domain/core/order/model"
	"github.com/D1sordxr/wb-tech-l0/internal/domain/core/order/ports"
)

type Cache struct {
	log         appPorts.Logger
	mu          sync.RWMutex
	store       map[string]*cacheItem
	ttl         time.Duration
	stopChan    chan struct{}
	initializer ports.CacheInitializer
}

const (
	ttl   = time.Minute * 5
	limit = 100
)

type cacheItem struct {
	order     *model.Order
	expiresAt time.Time
}

func NewCache(
	log appPorts.Logger,
	initializer ports.CacheInitializer,
) *Cache {
	cache := &Cache{
		log:         log,
		store:       make(map[string]*cacheItem),
		ttl:         ttl,
		stopChan:    make(chan struct{}),
		initializer: initializer,
	}

	return cache
}

func (c *Cache) Set(orderUID string, order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[orderUID] = &cacheItem{
		order:     order,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Get(orderUID string) *model.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.store[orderUID]
	if !exists {
		return nil
	}

	if time.Now().After(item.expiresAt) {
		return nil
	}

	return item.order
}

func (c *Cache) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.store {
		if now.After(item.expiresAt) {
			delete(c.store, key)
		}
	}
}

func (c *Cache) GetAll() map[string]*model.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]*model.Order)
	now := time.Now()

	for key, item := range c.store {
		if now.Before(item.expiresAt) {
			result[key] = item.order
		}
	}

	return result
}

func (c *Cache) Run(ctx context.Context) error {
	const op = "memory.Cache.Run"

	orders, err := c.initializer.GetOrdersForCache(ctx, limit)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("Successfully got orders for cache",
		"operation", op,
		"limit", limit,
		"orders_count", len(orders),
	)

	if len(orders) > 0 {
		lastUID := orders[len(orders)-1].OrderUID

		c.log.Info("Cache initialization",
			"operation", op,
			"orders_count", len(orders),
			"last_uid", lastUID,
		)

		for _, order := range orders {
			c.Set(order.OrderUID, order)
		}
	} else {
		c.log.Warn("No orders found for cache initialization")
	}

	cleanupTicker := time.NewTicker(c.ttl / 2)
	defer cleanupTicker.Stop()

	for {
		select {
		case <-cleanupTicker.C:
			c.cleanupExpired()
		case <-c.stopChan:
			return nil
		case <-ctx.Done():
			return nil
		}
	}
}
func (c *Cache) Shutdown(_ context.Context) error {
	close(c.stopChan)
	return nil
}
