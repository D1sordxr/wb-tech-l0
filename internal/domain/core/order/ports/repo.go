package ports

import (
	"context"
	"wb-tech-l0/internal/domain/core/order/model"
)

type OrderRepo interface {
	GetOrder(ctx context.Context, orderID string) (*model.Order, error)
	CreateOrder(ctx context.Context, order *model.Order) error
}

type CacheInitializer interface {
	GetOrdersForCache(ctx context.Context, limit int) ([]*model.Order, error)
}
