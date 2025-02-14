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
	c.App.POST("/api/v1/user/register")
	c.App.POST("/api/v1/user/admin/register")
	c.App.POST("/api/v1/user/login")
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.DELETE("/api/v1/user/logout")
	c.App.DELETE("/api/v1/user/delete")
}
