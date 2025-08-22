package reader

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"time"
	appPorts "wb-tech-l0/internal/domain/app/ports"
	"wb-tech-l0/internal/domain/core/order/ports"
	"wb-tech-l0/internal/infrastructure/kafka"
	"wb-tech-l0/internal/transport/kafka/order/dto"
)

type Reader struct {
	log       appPorts.Logger
	reader    *kafka.Reader
	validator *validator.Validate
	uc        ports.UseCase
	topic     string
}

func NewReader(
	log appPorts.Logger,
	reader *kafka.Reader,
	uc ports.UseCase,
	topic string,
) *Reader {
	return &Reader{
		log:       log,
		reader:    reader,
		uc:        uc,
		topic:     topic,
		validator: validator.New(),
	}
}

func (r *Reader) processMessage(ctx context.Context) {
	const op = "kafka.Reader.processMessage"
	withFields := func(args ...any) []any {
		return append([]any{"op", op}, args...)
	}

	message, err := r.reader.FetchMessage(ctx)
	if err != nil {
		r.log.Error("failed to fetch message", withFields("error", err.Error())...)
		return
	}

	var commitErr error
	defer func() {
		if commitErr == nil {
			if kafkaCommitErr := r.reader.CommitMessages(ctx, message); kafkaCommitErr != nil {
				r.log.Error("failed to commit message", withFields("error", kafkaCommitErr.Error())...)
			}
		} else {
			r.log.Error("failed to commit message due storage error", withFields("error", commitErr.Error())...)
		}
	}()

	if message.Topic != r.topic {
		r.log.Error("expected message to have exact topic",
			withFields("message_topic", message.Topic)...,
		)
		return
	}

	var msg dto.Order
	if err = json.Unmarshal(message.Value, &msg); err != nil {
		r.log.Error("failed to unmarshal message", withFields("error", err.Error())...)
		return
	}

	if err = r.validator.Struct(msg); err != nil {
		r.log.Error("failed to validate message", withFields("error", err.Error())...)
	}
	if commitErr = r.uc.CreateOrder(ctx, msg); err != nil {
		r.log.Error("failed to create order", withFields("error", err.Error())...)
		return
	}
}

func (r *Reader) Start(ctx context.Context) error {
	const op = "kafka.Reader.Start"
	withFields := func(args ...any) []any {
		return append([]any{"op", op}, args...)
	}

	r.log.Info("Starting kafka reader", withFields()...)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			r.processMessage(ctx)
		}
	}
}

func (r *Reader) Stop(_ context.Context) error {
	r.log.Info("Stopping kafka reader")
	return nil
}
