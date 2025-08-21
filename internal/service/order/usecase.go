package order

import (
	"context"
	"wb-tech-l0/internal/domain/ports"
)

type UseCase struct {
	log ports.Logger
	// repo or storage
}

func NewUseCase(log ports.Logger) *UseCase {
	return &UseCase{
		log: log,
	}
}

func (uc *UseCase) GetByID(ctx context.Context, orderID string) (string, error) {
	//TODO implement me
	panic("implement me")
}
