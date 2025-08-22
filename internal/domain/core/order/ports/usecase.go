package ports

import (
	"context"
	"wb-tech-l0/internal/domain/core/order/model"
	"wb-tech-l0/internal/transport/kafka/order/dto"
)

type UseCase interface {
	CreateOrder(ctx context.Context, orderDTO dto.Order) error
	GetByID(ctx context.Context, orderID string) (*model.Order, error)
}
