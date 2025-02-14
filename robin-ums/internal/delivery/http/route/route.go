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
	AuthMiddleware   gin.HandlerFunc
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.GET("/api/healthcheck", c.HealthController.Healthcheck)
	c.App.POST("/api/v1/user/register", c.UserController.RegisterUser)
	c.App.POST("/api/v1/user/admin/register", c.UserController.RegisterAdmin)
	c.App.POST("/api/v1/user/login", c.UserController.Login)
	c.App.GET("/api/v1/user/:id", c.UserController.GetUser)

	c.App.POST("/api/v1/artist/register", c.ArtistController.RegisterArtist)
	c.App.POST("/api/v1/artist/login", c.ArtistController.LoginArtist)
	c.App.GET("/api/v1/artist/:id", c.ArtistController.GetArtist)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.DELETE("/api/v1/user/logout", c.UserController.Logout)
	c.App.DELETE("/api/v1/user/delete/:id", c.UserController.DeleteUser)
	c.App.PUT("/api/v1/user/update", c.UserController.UpdateUser)

	c.App.PUT("/api/v1/artist/update", c.ArtistController.UpdateArtist)
	c.App.DELETE("/api/v1/artist/delete/:id", c.ArtistController.DeleteArtist)
	c.App.DELETE("/api/v1/artist/logout", c.ArtistController.LogoutArtist)
}

// TODO: exclusive admin route
// func (c *RouteConfig) SetupAdminRoute() {}
