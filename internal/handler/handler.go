package handler

import (
	"context"
	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Service interface {
	CheckImage(ctx context.Context, input domain.CheckImageRequest) (bool, error)
}

type Handler struct {
	log     *slog.Logger
	service Service
}

func NewHandler(log *slog.Logger, service Service) *Handler {
	return &Handler{log: log, service: service}
}

func (h *Handler) Check(c *gin.Context) {
	var input domain.CheckImageRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	exists, err := h.service.CheckImage(c.Request.Context(), input)

	if err != nil {
		slog.Error("error while check docker image", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": exists,
	})
}
