package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(r *gin.Engine, h *Handler) {
	r.GET("/check-image", h.Check)
}
