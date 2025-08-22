package order

import (
	"context"
	"errors"
	"fmt"

	orderErrs "github.com/D1sordxr/wb-tech-l0/internal/domain/core/order/errors"
	"github.com/D1sordxr/wb-tech-l0/internal/domain/core/order/model"
	"github.com/D1sordxr/wb-tech-l0/internal/infrastructure/storage/postgres"
	"github.com/D1sordxr/wb-tech-l0/internal/infrastructure/storage/postgres/repositories/order/gen"
	"github.com/D1sordxr/wb-tech-l0/internal/infrastructure/storage/postgres/tools"

	"github.com/jackc/pgx/v5"
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

func (r *Repository) GetOrder(ctx context.Context, orderUID string) (*model.Order, error) {
	const op = "repositories.order.GetOrder"

	tx, err := r.executor.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin tx: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := r.queries.WithTx(tx)

	orderDB, err := qtx.GetOrder(ctx, orderUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, orderErrs.ErrOrderNotFount
		}
		return nil, fmt.Errorf("%s: failed to get order: %w", op, err)
	}

	deliveryDB, err := qtx.GetDelivery(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get delivery: %w", op, err)
	}

	paymentDB, err := qtx.GetPayment(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get payment: %w", op, err)
	}

	itemsDB, err := qtx.GetItems(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get items: %w", op, err)
	}

	items := make([]model.Item, len(itemsDB))
	for i, item := range itemsDB {
		items[i] = model.Item{
			ChrtID:      item.ChrtID.Int64,
			TrackNumber: item.TrackNumber.String,
			Price:       item.Price.Int32,
			RID:         item.Rid.String,
			Name:        item.ItemName.String,
			Sale:        item.Sale.Int32,
			Size:        item.ItemSize.String,
			TotalPrice:  item.TotalPrice.Int32,
			NmID:        item.NmID.Int64,
			Brand:       item.Brand.String,
			Status:      item.Status.Int32,
		}
	}

	payment := model.Payment{
		Transaction:  paymentDB.TransactionID,
		RequestID:    paymentDB.RequestID.String,
		Currency:     paymentDB.Currency.String,
		Provider:     paymentDB.Provider.String,
		Amount:       paymentDB.Amount.Int32,
		PaymentDt:    paymentDB.PaymentDt.Int64,
		Bank:         paymentDB.Bank.String,
		DeliveryCost: paymentDB.DeliveryCost.Int32,
		GoodsTotal:   paymentDB.GoodsTotal.Int32,
		CustomFee:    paymentDB.CustomFee.Int32,
	}

	delivery := model.Delivery{
		Name:    deliveryDB.DelName,
		Phone:   deliveryDB.Phone,
		Zip:     deliveryDB.Zip.String,
		City:    deliveryDB.City.String,
		Address: deliveryDB.Address.String,
		Region:  deliveryDB.Region.String,
		Email:   deliveryDB.Email.String,
	}

	order := &model.Order{
		OrderUID:          orderDB.OrderUid,
		TrackNumber:       orderDB.TrackNumber,
		Entry:             orderDB.Entry,
		Locale:            orderDB.Locale,
		InternalSignature: orderDB.InternalSignature.String,
		CustomerID:        orderDB.CustomerID,
		DeliveryService:   orderDB.DeliveryService.String,
		ShardKey:          orderDB.Shardkey.String,
		SmID:              orderDB.SmID,
		DateCreated:       orderDB.DateCreated.Time,
		OofShard:          orderDB.OofShard.String,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit tx: %w", op, err)
	}

	return order, nil
}

func (r *Repository) CreateOrder(ctx context.Context, order *model.Order) error {
	const op = "repositories.order.CreateOrder"

	tx, err := r.executor.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin tx: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := r.queries.WithTx(tx)

	err = qtx.CreateOrder(ctx, gen.CreateOrderParams{
		OrderUid:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Locale:            order.Locale,
		InternalSignature: tools.ToText(order.InternalSignature),
		CustomerID:        order.CustomerID,
		DeliveryService:   tools.ToText(order.DeliveryService),
		Shardkey:          tools.ToText(order.ShardKey),
		SmID:              order.SmID,
		DateCreated:       tools.ToTimestamp(order.DateCreated),
		OofShard:          tools.ToText(order.OofShard),
	})
	if err != nil {
		if tools.IsUniqueErr(err) {
			return orderErrs.ErrOrderAlreadyExists
		}
		return fmt.Errorf("%s: faield to create order: %w", op, err)
	}

	itemsRows := make([][]any, len(order.Items))
	for i, item := range order.Items {
		itemsRows[i] = []any{
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		}
	}
	itemColumnNames := []string{
		"order_uid",
		"chrt_id",
		"track_number",
		"price",
		"rid",
		"item_name",
		"sale",
		"item_size",
		"total_price",
		"nm_id",
		"brand",
		"status",
	}
	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"items"},
		itemColumnNames,
		pgx.CopyFromRows(itemsRows),
	)
	if err != nil {
		return fmt.Errorf("%s: failed to insert items: %w", op, err)
	}

	err = qtx.CreateDelivery(ctx, gen.CreateDeliveryParams{
		OrderUid: order.OrderUID,
		DelName:  order.Delivery.Name,
		Phone:    order.Delivery.Phone,
		Zip:      tools.ToText(order.Delivery.Zip),
		City:     tools.ToText(order.Delivery.City),
		Address:  tools.ToText(order.Delivery.Address),
		Region:   tools.ToText(order.Delivery.Region),
		Email:    tools.ToText(order.Delivery.Email),
	})
	if err != nil {
		return fmt.Errorf("%s: failed to create delivery: %w", op, err)
	}

	err = qtx.CreatePayment(ctx, gen.CreatePaymentParams{
		OrderUid:      order.OrderUID,
		TransactionID: order.Payment.Transaction,
		RequestID:     tools.ToText(order.Payment.RequestID),
		Currency:      tools.ToText(order.Payment.Currency),
		Provider:      tools.ToText(order.Payment.Provider),
		Amount:        tools.ToInt4(order.Payment.Amount),
		PaymentDt:     tools.ToInt8(order.Payment.PaymentDt),
		Bank:          tools.ToText(order.Payment.Bank),
		DeliveryCost:  tools.ToInt4(order.Payment.DeliveryCost),
		GoodsTotal:    tools.ToInt4(order.Payment.GoodsTotal),
		CustomFee:     tools.ToInt4(order.Payment.CustomFee),
	})
	if err != nil {
		return fmt.Errorf("%s: failed to create payment: %w", op, err)
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetOrdersForCache(ctx context.Context, limit int) ([]*model.Order, error) {
	const op = "repositories.order.GetOrdersForCache"

	tx, err := r.executor.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin tx: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := r.queries.WithTx(tx)

	ordersDB, err := qtx.GetLatestOrders(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get orders: %w", op, err)
	}

	if len(ordersDB) == 0 {
		return []*model.Order{}, nil
	}

	orderIDs := make([]string, len(ordersDB))
	for i, order := range ordersDB {
		orderIDs[i] = order.OrderUid
	}

	items, err := qtx.GetItemsForOrders(ctx, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get items: %w", op, err)
	}

	deliveries, err := qtx.GetDeliveriesForOrders(ctx, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get deliveries: %w", op, err)
	}

	payments, err := qtx.GetPaymentsForOrders(ctx, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get payments: %w", op, err)
	}

	itemsMap := make(map[string][]gen.Item)
	for _, item := range items {
		itemsMap[item.OrderUid] = append(itemsMap[item.OrderUid], item)
	}

	deliveriesMap := make(map[string]gen.Delivery)
	for _, delivery := range deliveries {
		deliveriesMap[delivery.OrderUid] = delivery
	}

	paymentsMap := make(map[string]gen.Payment)
	for _, payment := range payments {
		paymentsMap[payment.OrderUid] = payment
	}

	orders := make([]*model.Order, 0, len(ordersDB))

	for _, orderDB := range ordersDB {
		orderUID := orderDB.OrderUid

		delivery, deliveryExists := deliveriesMap[orderUID]
		payment, paymentExists := paymentsMap[orderUID]
		orderItems, itemsExist := itemsMap[orderUID]

		if !deliveryExists || !paymentExists || !itemsExist {
			continue
		}

		modelItems := make([]model.Item, len(orderItems))
		for i, item := range orderItems {
			modelItems[i] = model.Item{
				ChrtID:      item.ChrtID.Int64,
				TrackNumber: item.TrackNumber.String,
				Price:       item.Price.Int32,
				RID:         item.Rid.String,
				Name:        item.ItemName.String,
				Sale:        item.Sale.Int32,
				Size:        item.ItemSize.String,
				TotalPrice:  item.TotalPrice.Int32,
				NmID:        item.NmID.Int64,
				Brand:       item.Brand.String,
				Status:      item.Status.Int32,
			}
		}

		modelDelivery := model.Delivery{
			Name:    delivery.DelName,
			Phone:   delivery.Phone,
			Zip:     delivery.Zip.String,
			City:    delivery.City.String,
			Address: delivery.Address.String,
			Region:  delivery.Region.String,
			Email:   delivery.Email.String,
		}

		modelPayment := model.Payment{
			Transaction:  payment.TransactionID,
			RequestID:    payment.RequestID.String,
			Currency:     payment.Currency.String,
			Provider:     payment.Provider.String,
			Amount:       payment.Amount.Int32,
			PaymentDt:    payment.PaymentDt.Int64,
			Bank:         payment.Bank.String,
			DeliveryCost: payment.DeliveryCost.Int32,
			GoodsTotal:   payment.GoodsTotal.Int32,
			CustomFee:    payment.CustomFee.Int32,
		}

		order := &model.Order{
			OrderUID:          orderDB.OrderUid,
			TrackNumber:       orderDB.TrackNumber,
			Entry:             orderDB.Entry,
			Locale:            orderDB.Locale,
			InternalSignature: orderDB.InternalSignature.String,
			CustomerID:        orderDB.CustomerID,
			DeliveryService:   orderDB.DeliveryService.String,
			ShardKey:          orderDB.Shardkey.String,
			SmID:              orderDB.SmID,
			DateCreated:       orderDB.DateCreated.Time,
			OofShard:          orderDB.OofShard.String,
			Delivery:          modelDelivery,
			Payment:           modelPayment,
			Items:             modelItems,
		}

		orders = append(orders, order)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit tx: %w", op, err)
	}

	return orders, nil
}
