package job

import (
	"context"
	"encoding/json"
	kafkaLib "github.com/segmentio/kafka-go"
	"time"
	appPorts "wb-tech-l0/internal/domain/app/ports"
	"wb-tech-l0/internal/infrastructure/kafka"
	"wb-tech-l0/internal/infrastructure/mock"
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
