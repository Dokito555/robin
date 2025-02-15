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
	userGroup := c.App.Group("/api/v1/user")
	userGroup.Use(c.UserMiddleware)
	{
		userGroup.DELETE("/api/v1/user/logout", c.UserController.Logout)
		userGroup.PUT("/api/v1/user/update", c.UserController.UpdateUser)
	}
}

func (c *RouteConfig) SetupArtistRoute() {
	artistGroup := c.App.Group("/api/v1/artist")
	artistGroup.Use(c.ArtistMiddelware)
	{
		artistGroup.PUT("/update", c.ArtistController.UpdateArtist)
		artistGroup.DELETE("/logout", c.ArtistController.LogoutArtist)
		artistGroup.POST("/register", c.ArtistController.RegisterArtist)
		artistGroup.POST("/login", c.ArtistController.LoginArtist)
		artistGroup.GET("/:id", c.ArtistController.GetArtist)
	}
}

func (c *RouteConfig) SetupAdminRoute() {
	adminGroup := c.App.Group("/api/v1/user/admin")
	adminGroup.Use(c.AdminMiddleware)
	{
		adminGroup.POST("/register", c.UserController.RegisterAdmin)
		adminGroup.DELETE("/artist/delete/:id", c.ArtistController.DeleteArtist)
		adminGroup.DELETE("/user/delete/:id", c.UserController.DeleteUser)
	}
}
