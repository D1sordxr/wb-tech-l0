package mapper

import (
	"github.com/D1sordxr/wb-tech-l0/internal/domain/core/order/model"
	"github.com/D1sordxr/wb-tech-l0/internal/transport/kafka/order/dto"
)

func OrderFromDTO(dtoOrder dto.Order) *model.Order {
	return &model.Order{
		OrderUID:          dtoOrder.ID,
		TrackNumber:       dtoOrder.TrackNumber,
		Entry:             dtoOrder.Entry,
		Locale:            dtoOrder.Locale,
		InternalSignature: dtoOrder.InternalSignature,
		CustomerID:        dtoOrder.CustomerID,
		DeliveryService:   dtoOrder.DeliveryService,
		ShardKey:          dtoOrder.ShardKey,
		SmID:              dtoOrder.SmID,
		DateCreated:       dtoOrder.DateCreated,
		OofShard:          dtoOrder.OofShard,
		Delivery:          deliveryFromDTO(dtoOrder.Delivery),
		Payment:           paymentFromDTO(dtoOrder.Payment),
		Items:             itemsFromDTO(dtoOrder.Items),
	}
}

func deliveryFromDTO(dtoDelivery dto.Delivery) model.Delivery {
	return model.Delivery{
		Name:    dtoDelivery.Name,
		Phone:   dtoDelivery.Phone,
		Zip:     dtoDelivery.Zip,
		City:    dtoDelivery.City,
		Address: dtoDelivery.Address,
		Region:  dtoDelivery.Region,
		Email:   dtoDelivery.Email,
	}
}

func paymentFromDTO(dtoPayment dto.Payment) model.Payment {
	return model.Payment{
		Transaction:  dtoPayment.Transaction,
		RequestID:    dtoPayment.RequestID,
		Currency:     dtoPayment.Currency,
		Provider:     dtoPayment.Provider,
		Amount:       dtoPayment.Amount,
		PaymentDt:    dtoPayment.PaymentDt,
		Bank:         dtoPayment.Bank,
		DeliveryCost: dtoPayment.DeliveryCost,
		GoodsTotal:   dtoPayment.GoodsTotal,
		CustomFee:    dtoPayment.CustomFee,
	}
}

func itemsFromDTO(dtoItems []dto.Item) []model.Item {
	items := make([]model.Item, len(dtoItems))
	for i, dtoItem := range dtoItems {
		items[i] = model.Item{
			ChrtID:      dtoItem.ChrtID,
			TrackNumber: dtoItem.TrackNumber,
			Price:       dtoItem.Price,
			RID:         dtoItem.RID,
			Name:        dtoItem.Name,
			Sale:        dtoItem.Sale,
			Size:        dtoItem.Size,
			TotalPrice:  dtoItem.TotalPrice,
			NmID:        dtoItem.NmID,
			Brand:       dtoItem.Brand,
			Status:      dtoItem.Status,
		}
	}
	return items
}
