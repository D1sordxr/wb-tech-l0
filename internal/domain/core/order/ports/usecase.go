package ports

import "context"

type UseCase interface {
	GetByID(ctx context.Context, orderID string) (string, error)
}
