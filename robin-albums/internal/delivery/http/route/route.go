package route

import (
	"github.com/Dokito555/robin-albums/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App              *gin.Engine
	HealthController *http.HealthController
	AlbumController  *http.AlbumController
	AuthMiddleware   gin.HandlerFunc
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.GET("/api/healthcheck", c.HealthController.Healthcheck)
	c.App.GET("/api/v1/album/:id", c.AlbumController.GetAlbum)
}

func (c *RouteConfig) SetupAuthRoute() {
	albumGroup := c.App.Group("/api/v1/album")
	albumGroup.Use(c.AuthMiddleware)
	{
		albumGroup.POST("", c.AlbumController.CreateAlbum)
		albumGroup.PUT("/update", c.AlbumController.UpdateAlbum)
		albumGroup.DELETE("/:id", c.AlbumController.DeleteAlbum)
	}
}
