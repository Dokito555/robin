package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *Handler) Ready(c *gin.Context) {
	res := h.service.Ready()

	if !res.Ok {
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}

	c.JSON(http.StatusOK, res)
}