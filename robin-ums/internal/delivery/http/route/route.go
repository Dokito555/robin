package route

import (
	"github.com/Dokito555/robin-ums/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App              *gin.Engine
	HealthController *http.HealthController
	UserController   *http.UserController
	ArtistController *http.ArtistController
	UserMiddleware   gin.HandlerFunc
	AdminMiddleware  gin.HandlerFunc
	ArtistMiddelware gin.HandlerFunc
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
	c.SetupAdminRoute()
	c.SetupArtistRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.GET("/api/healthcheck", c.HealthController.Healthcheck)
	c.App.POST("/api/v1/user/register", c.UserController.RegisterUser)
	c.App.POST("/api/v1/user/login", c.UserController.Login)
	c.App.GET("/api/v1/user/:id", c.UserController.GetUser)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.UserMiddleware)
	c.App.DELETE("/api/v1/user/logout", c.UserController.Logout)
	c.App.PUT("/api/v1/user/update", c.UserController.UpdateUser)
}

func (c *RouteConfig) SetupArtistRoute() {
	artistGroup := c.App.Group("/api/v1/artist")
	artistGroup.Use(c.ArtistMiddelware)
	{
		c.App.PUT("/api/v1/artist/update", c.ArtistController.UpdateArtist)
		c.App.DELETE("/api/v1/artist/logout", c.ArtistController.LogoutArtist)
		c.App.POST("/api/v1/artist/register", c.ArtistController.RegisterArtist)
		c.App.POST("/api/v1/artist/login", c.ArtistController.LoginArtist)
		c.App.GET("/api/v1/artist/:id", c.ArtistController.GetArtist)
	}
}

func (c *RouteConfig) SetupAdminRoute() {
	adminGroup := c.App.Group("/api/v1/user/admin")
	adminGroup.Use(c.AdminMiddleware)
	{
		c.App.POST("/api/v1/user/admin/register", c.UserController.RegisterAdmin)
		c.App.DELETE("/api/v1/artist/delete/:id", c.ArtistController.DeleteArtist)
		c.App.DELETE("/api/v1/user/delete/:id", c.UserController.DeleteUser)
	}
}
