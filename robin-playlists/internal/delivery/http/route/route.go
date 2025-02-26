package route

import (
	"github.com/Dokito555/robin-playlists/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App                *gin.Engine
	HealthController   *http.HealthController
	PlaylistController *http.PlaylistController
	AuthMiddleware     gin.HandlerFunc
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.GET("/api/healthcheck", c.HealthController.Healthcheck)
	c.App.GET("/api/v1/playlist/:id", c.PlaylistController.GetPlaylist)
}

func (c *RouteConfig) SetupAuthRoute() {
	playlistGroup := c.App.Group("/api/v1/playlist")
	playlistGroup.Use(c.AuthMiddleware)
	{
		playlistGroup.POST("", c.PlaylistController.CreatePlaylist)
		playlistGroup.PUT("/:id", c.PlaylistController.UpdatePlaylist)
		playlistGroup.DELETE("/:id", c.PlaylistController.DeletePlaylist)
		playlistGroup.GET("", c.PlaylistController.GetPlaylistList)
	}
}
