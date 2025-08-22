package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"wb-tech-l0/internal/app"
	"wb-tech-l0/internal/config"
	"wb-tech-l0/internal/service/order"
	"wb-tech-l0/internal/storage/postgres"
	orderRepopository "wb-tech-l0/internal/storage/postgres/repositories/order"
	"wb-tech-l0/internal/transport/http"
	"wb-tech-l0/internal/transport/http/order/handler"
	"wb-tech-l0/internal/transport/kafka"
	"wb-tech-l0/internal/transport/kafka/reader"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()

	log := slog.Default()

	pool := postgres.NewPool(ctx, &cfg.Storage)

	orderRepo := orderRepopository.NewOrderRepo(pool)

	orderUseCase := order.NewUseCase(log, orderRepo) // todo

	orderHandler := handler.NewHandler(orderUseCase) // todo

	httpServer := http.NewServer(
		log,
		&cfg.Server,
		orderHandler,
	)

	orderReader := reader.NewReader() // todo

	kafkaWorker := kafka.NewWorker( // todo
		log,
		orderReader,
	)

	appContainer := app.NewApp(
		log,
		pool,
		httpServer,
		kafkaWorker,
	)
	appContainer.Run(ctx)
}
