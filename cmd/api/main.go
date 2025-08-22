package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"wb-tech-l0/internal/infrastructure/app"
	cache "wb-tech-l0/internal/infrastructure/cache/memory/order"
	"wb-tech-l0/internal/infrastructure/config"
	"wb-tech-l0/internal/infrastructure/kafka"
	"wb-tech-l0/internal/infrastructure/storage/postgres"
	orderRepopository "wb-tech-l0/internal/infrastructure/storage/postgres/repositories/order"
	loadWorker "wb-tech-l0/internal/infrastructure/worker"
	"wb-tech-l0/internal/infrastructure/worker/job"
	"wb-tech-l0/internal/service/order"
	"wb-tech-l0/internal/transport/http"
	"wb-tech-l0/internal/transport/http/order/handler"
	"wb-tech-l0/internal/transport/kafka/order/reader"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()

	log := slog.Default()

	pool := postgres.NewPool(ctx, &cfg.Storage)
	orderRepo := orderRepopository.NewOrderRepo(pool)

	orderCache := cache.NewCache(log, orderRepo)

	orderReaderConn := kafka.NewReader(
		log,
		&cfg.MessageBroker,
		cfg.MessageBroker.SaverGroup,
	)
	orderWriterConn := kafka.NewWriter(log, &cfg.MessageBroker)

	orderUseCase := order.NewUseCase(
		log,
		orderRepo,
		orderCache,
	)

	orderHandler := handler.NewHandler(orderUseCase)

	httpServer := http.NewServer(
		log,
		&cfg.Server,
		orderHandler,
	)

	orderKafkaReader := reader.NewReader(
		log,
		orderReaderConn,
		orderUseCase,
		cfg.MessageBroker.OrdersTopic,
	)

	orderKafkaWriter := job.NewMockOrderWriter(log, orderWriterConn)

	worker := loadWorker.NewWorker(
		log,
		orderKafkaReader,
		orderKafkaWriter,
	)

	appContainer := app.NewApp(
		log,
		pool,
		orderCache,
		orderReaderConn,
		orderWriterConn,
		httpServer,
		worker,
	)
	appContainer.Run(ctx)
}
