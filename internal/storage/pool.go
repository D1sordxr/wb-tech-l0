package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"wb-tech-l0/internal/config"
)

type Pool struct {
	*pgxpool.Pool
}

func NewPool(ctx context.Context, config *config.Postgres) *Pool {
	pool, err := pgxpool.New(ctx, config.ConnectionString())
	if err != nil {
		panic(err)
	}

	return &Pool{Pool: pool}
}

func (p *Pool) Run(ctx context.Context) error {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			return p.Pool.Ping(ctx)
		}
	}
}

func (p *Pool) Shutdown(_ context.Context) error {
	p.Pool.Close()
	return nil
}
