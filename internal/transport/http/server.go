package http

import (
	"context"
	"errors"
	"github.com/D1sordxr/wb-tech-l0/internal/domain/app/ports"
	"github.com/D1sordxr/wb-tech-l0/internal/infrastructure/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Handler interface {
	RegisterRoutes(router gin.IRouter)
}

type Server struct {
	log      ports.Logger
	handlers []Handler
	engine   *gin.Engine
	server   *http.Server
}

func NewServer(
	log ports.Logger,
	config *config.HTTPServer,
	handlers ...Handler,
) *Server {
	log.Info("Initializing HTTP server", "port", config.Port)

	engine := gin.Default()

	if config.CORS {
		allowedOrigins := config.AllowOrigins
		if len(allowedOrigins) == 0 {
			allowedOrigins = []string{"*"}
		}

		engine.Use(cors.New(cors.Config{
			AllowOrigins:     allowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	return &Server{
		log: log,
		server: &http.Server{
			Addr:              ":" + config.Port,
			Handler:           engine.Handler(),
			ReadHeaderTimeout: config.Timeout,
			ReadTimeout:       config.Timeout,
			WriteTimeout:      config.Timeout,
		},
		engine:   engine,
		handlers: handlers,
	}
}

func (s *Server) Run(_ context.Context) error {
	s.log.Info("Registering HTTP handlers...")
	for _, handler := range s.handlers {
		group := s.engine.Group("/api")
		handler.RegisterRoutes(group)
	}

	s.log.Info("Starting HTTP server...", "address", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("HTTP server closed gracefully")
			return nil
		}
		s.log.Error("HTTP server stopped with error", "error", err.Error())
		return err
	}

	s.log.Info("HTTP server exited unexpectedly")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down HTTP server...")
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("Failed to gracefully shutdown HTTP server", "error", err.Error())
		return err
	}
	s.log.Info("HTTP server shutdown complete")
	return nil
}
