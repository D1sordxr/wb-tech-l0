package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/D1sordxr/wb-tech-l0/internal/domain/core/order/ports"
	sharedErrs "github.com/D1sordxr/wb-tech-l0/internal/domain/core/shared/errors"
	"github.com/D1sordxr/wb-tech-l0/pkg/errtool"

	"github.com/gin-gonic/gin"
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
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "id required"})
	}

	resp, err := h.getOrderUseCase.GetByID(reqCtx, id)
	if err != nil {
		switch {
		case errtool.In(
			err,
			sharedErrs.ErrOrderUIDInvalidChars,
			sharedErrs.ErrOrderUIDInvalidSuffix,
			sharedErrs.ErrOrderUIDInvalidLength,
		):
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, &resp)
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/order/:id", h.getByID)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
}
