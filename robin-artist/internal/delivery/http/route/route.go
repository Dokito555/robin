package route

import (
	"github.com/Dokito555/robin-artist/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App              *gin.Engine
	HealthController *http.HealthController
	ArtistController *http.ArtistController
	AuthMiddleware   gin.HandlerFunc
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.GET("/api/healthcheck", c.HealthController.Healthcheck)
	c.App.POST("/api/v1/artist/register", c.ArtistController.RegisterArtist)
	c.App.GET("/api/v1/artist/:id", c.ArtistController.GetArtist)
	c.App.POST("/api/v1/artist/login", c.ArtistController.LoginArtist)
}

func (c *RouteConfig) SetupAuthRoute() {
	artistGroup := c.App.Group("/api/v1/artist")
	artistGroup.Use(c.AuthMiddleware)
	{
		artistGroup.PUT("/update", c.ArtistController.UpdateArtist)
		artistGroup.DELETE("/logout", c.ArtistController.LogoutArtist)
	}
}
