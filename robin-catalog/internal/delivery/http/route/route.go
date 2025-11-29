package route

import (
	"github.com/Dokito555/robin/robin-catalog/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App              *gin.Engine
	HealthController *http.HealthController
	ArtistController *http.ArtistController
	AlbumController  *http.AlbumController
	// SongController   *http.SongController
	AuthMiddleware   gin.HandlerFunc
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.GET("/api/healthcheck", c.HealthController.Healthcheck)

	// Artist
	c.App.POST("/api/v1/artist/register", c.ArtistController.RegisterArtist)
	c.App.GET("/api/v1/artist/:id", c.ArtistController.GetArtist)
	c.App.POST("/api/v1/artist/login", c.ArtistController.LoginArtist)
	c.App.GET("/api/v1/artist/list", c.ArtistController.GetArtistList)

	// Album
	c.App.GET("/api/v1/album/:id", c.AlbumController.GetAlbum)
	c.App.GET("/api/v1/album/artist/:id", c.AlbumController.GetAlbumListByArtistID)
}

func (c *RouteConfig) SetupAuthRoute() {
	// Artist
	artistGroup := c.App.Group("/api/v1/artist")
	artistGroup.Use(c.AuthMiddleware)
	{
		artistGroup.PUT("/update", c.ArtistController.UpdateArtist)
		artistGroup.DELETE("/logout", c.ArtistController.LogoutArtist)
	}

	// Album
	albumGroup := c.App.Group("/api/v1/album")
	albumGroup.Use(c.AuthMiddleware)
	{
		albumGroup.POST("", c.AlbumController.CreateAlbum)
		albumGroup.PUT("/:id", c.AlbumController.UpdateAlbum)
		albumGroup.DELETE("/:id", c.AlbumController.DeleteAlbum)
	}
}
