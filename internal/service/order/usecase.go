package order

import (
	"context"
	"fmt"
	appPorts "wb-tech-l0/internal/domain/app/ports"
	"wb-tech-l0/internal/domain/core/order/ports"
	"wb-tech-l0/internal/domain/core/shared/vo"
	"wb-tech-l0/internal/storage/postgres/repositories/order/gen"
)

type UseCase struct {
	log  appPorts.Logger
	repo ports.OrderRepo
}

func NewUseCase(
	log appPorts.Logger,
	repo ports.OrderRepo,
) *UseCase {
	return &UseCase{
		log:  log,
		repo: repo,
	}
}

func (uc *UseCase) GetByID(
	ctx context.Context,
	orderID string,
) (
	*gen.Order,
	error,
) {
	const op = "service.order.UseCase.GetByID"
	withFields := func(args ...any) []any {
		return append([]any{"op", op, "orderUID", orderID}, args...)
	}

	uc.log.Info("Attempting to get order", withFields()...)

	if err := vo.ValidateUID(orderID); err != nil {
		uc.log.Error("Failed to validate order", withFields("error", err.Error())...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	orderModel, err := uc.repo.GetOrder(ctx, orderID)
	if err != nil {
		uc.log.Error("Failed to get order", withFields("error", err.Error())...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO parse and return dto

	uc.log.Info("Successfully got order", withFields()...)

	return orderModel, nil
}
