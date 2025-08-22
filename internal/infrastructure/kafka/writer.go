package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
	"wb-tech-l0/internal/domain/app/ports"
	"wb-tech-l0/internal/infrastructure/config"
)

type Writer struct {
	log ports.Logger
	*kafka.Writer
	address       string
	topic         string
	isCreateTopic bool
}

func NewWriter(log ports.Logger, cfg *config.Kafka) *Writer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP([]string{cfg.Address}...),
		Balancer: &kafka.LeastBytes{},
	}

	return &Writer{
		log:           log,
		address:       cfg.Address,
		topic:         cfg.OrdersTopic,
		Writer:        writer,
		isCreateTopic: cfg.CreateTopic,
	}
}

func (w *Writer) GetTopic() string {
	return w.topic
}

const (
	partitions        = 3
	replicationFactor = 1
)

func (w *Writer) createTopic() error {
	conn, err := kafka.Dial("tcp", w.address)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	if err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             w.topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}); err != nil {
		return err
	}

	return nil
}

func (w *Writer) Run(ctx context.Context) error {
	const op = "kafka.Writer.Run"

	if w.isCreateTopic {
		if err := w.createTopic(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			w.log.Info("Kafka writer stats", "stats", w.Writer.Stats())
		}
	}
}

func (w *Writer) Shutdown(_ context.Context) error {
	_ = w.Close()
	return nil
}
