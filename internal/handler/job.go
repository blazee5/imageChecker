package handler

import (
	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) CreateJob(c *gin.Context) {
	var input domain.CreateJobRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	err := h.service.Job.CreateJob(c.Request.Context(), input)

	if err != nil {
		slog.Error("error while create job", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
