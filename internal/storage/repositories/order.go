package repositories

import "wb-tech-l0/internal/storage"

type OrderRepo struct {
	executor *storage.Pool
}

func NewOrderRepo(executor *storage.Pool) *OrderRepo {
	return &OrderRepo{executor: executor}
}
