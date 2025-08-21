package handler

import (
	"github.com/gin-gonic/gin"
	"wb-tech-l0/internal/domain/core/order/ports"
)

type Handler struct {
	getOrderUseCase ports.UseCase
}

func NewHandler(getOrderUseCase ports.UseCase) *Handler {
	return &Handler{
		getOrderUseCase: getOrderUseCase,
	}
}

func (h *Handler) getByID(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "not implemented",
	})
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/order/:id", h.getByID)
}
