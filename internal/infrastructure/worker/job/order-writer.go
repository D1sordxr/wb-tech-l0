package job

import (
	"context"
	"encoding/json"
	"time"

	appPorts "github.com/D1sordxr/wb-tech-l0/internal/domain/app/ports"
	"github.com/D1sordxr/wb-tech-l0/internal/infrastructure/kafka"
	"github.com/D1sordxr/wb-tech-l0/internal/infrastructure/mock"

	kafkaLib "github.com/segmentio/kafka-go"
)

type MockOrderWriter struct {
	log      appPorts.Logger
	orderGen *mock.Generator
	writer   *kafka.Writer
}

func NewMockOrderWriter(
	log appPorts.Logger,
	writer *kafka.Writer,
) *MockOrderWriter {
	return &MockOrderWriter{
		log:      log,
		orderGen: mock.NewMockGenerator(),
		writer:   writer,
	}
}

func (w *MockOrderWriter) Start(ctx context.Context) error {
	const op = "job.MockOrderWriter.Start"

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			order := w.orderGen.GenerateOrder()
			data, _ := json.Marshal(order)
			if err := w.writer.WriteMessages(ctx, kafkaLib.Message{
				Topic: w.writer.GetTopic(),
				Value: data,
			}); err != nil {
				w.log.Error("failed to write mock order", "op", op, "error", err.Error())
			}
		}
	}
}

func (w *MockOrderWriter) Stop(_ context.Context) error {
	return nil
}
