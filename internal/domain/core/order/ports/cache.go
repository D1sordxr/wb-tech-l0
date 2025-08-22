package ports

import "github.com/D1sordxr/wb-tech-l0/internal/domain/core/order/model"

type OrderCache interface {
	Set(orderUID string, order *model.Order)
	Get(orderUID string) *model.Order
}
