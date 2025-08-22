package kafka

import (
	"context"
	"time"

	appPorts "github.com/D1sordxr/wb-tech-l0/internal/domain/app/ports"
	"github.com/D1sordxr/wb-tech-l0/internal/infrastructure/config"

	"github.com/segmentio/kafka-go"
)

type Reader struct {
	log appPorts.Logger
	*kafka.Reader
}

func NewReader(log appPorts.Logger, cfg *config.Kafka, group string) *Reader {
	return &Reader{
		log: log,
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{cfg.Address},
			GroupTopics: []string{cfg.OrdersTopic},
			GroupID:     group,
		}),
	}
}

func (r *Reader) Run(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			r.log.Info("Kafka reader stats", "stats", r.Stats())
		}
	}
}

func (r *Reader) Shutdown(_ context.Context) error {
	_ = r.Close()
	return nil
}
