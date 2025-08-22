package ports

import (
	"context"
	"wb-tech-l0/internal/storage/postgres/repositories/order/gen"
)

type OrderRepo interface {
	GetOrder(ctx context.Context, orderID string) (*gen.Order, error)
}
