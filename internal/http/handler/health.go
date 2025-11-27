package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Health check
// @Description Check application health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *Handler) CheckHealth(c *gin.Context) {
	op := "http.handler.health"
	slog.Debug("Health check requested", "operation", op)

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
	})
}
