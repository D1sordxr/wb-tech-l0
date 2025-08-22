package order

import (
	"context"
	"fmt"
	"wb-tech-l0/internal/domain/core/order/model"
	"wb-tech-l0/internal/storage/postgres"
	"wb-tech-l0/internal/storage/postgres/repositories/order/gen"
	"wb-tech-l0/internal/storage/postgres/tools"
)

type Repository struct {
	executor *postgres.Pool
	queries  *gen.Queries
}

func NewOrderRepo(executor *postgres.Pool) *Repository {
	return &Repository{
		executor: executor,
		queries:  gen.New(executor),
	}
}

func (r *Repository) GetOrder(ctx context.Context, orderID string) (*gen.Order, error) {
	const op = "repositories.order.GetOrder"

	order, err := r.queries.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &order, nil
}

func (r *Repository) CreateOrder(ctx context.Context, orderData model.Order) error {
	const op = "repositories.order.CreateOrder"

	tx, err := r.executor.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := r.queries.WithTx(tx)

	if err = qtx.CreateOrder(ctx, gen.CreateOrderParams{
		OrderUid:          orderData.OrderUID,
		TrackNumber:       orderData.TrackNumber,
		Entry:             orderData.Entry,
		Locale:            orderData.Locale,
		InternalSignature: tools.ToText(orderData.InternalSignature),
		CustomerID:        orderData.CustomerID,
		DeliveryService:   tools.ToText(orderData.DeliveryService),
		Shardkey:          tools.ToText(orderData.ShardKey),
		SmID:              orderData.SmID,
		DateCreated:       tools.ToTimestamp(orderData.DateCreated),
		OofShard:          tools.ToText(orderData.OofShard),
	}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// if err := qtx.createSomethingIfNeeded(ctx, ...); err != nil { ... }

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit failed: %w", op, err)
	}

	return nil
}
