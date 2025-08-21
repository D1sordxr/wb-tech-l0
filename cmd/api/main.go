package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"wb-tech-l0/internal/app"
	"wb-tech-l0/internal/config"
	"wb-tech-l0/internal/service/order"
	"wb-tech-l0/internal/storage"
	"wb-tech-l0/internal/transport/http"
	"wb-tech-l0/internal/transport/http/order/handler"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()

	log := slog.Default()

	pool := storage.NewPool(ctx, &cfg.Storage)

	orderUseCase := order.NewUseCase(log) // todo

	orderHandler := handler.NewHandler(orderUseCase) // todo

	httpServer := http.NewServer(
		log,
		&cfg.Server,
		orderHandler,
	)

	appContainer := app.NewApp(
		log,
		pool,
		httpServer,
	)
	appContainer.Run(ctx)
}
