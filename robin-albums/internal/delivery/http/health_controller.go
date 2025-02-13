package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// healthcheck handlers
type HealthController struct {
	Log *logrus.Logger	
}

func NewHealthController(logger *logrus.Logger) *HealthController {
	return &HealthController{
		Log: logger,
	}
}

func (c *HealthController) Healthcheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "healthy",
	})
}