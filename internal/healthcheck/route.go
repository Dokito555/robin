package healthcheck

import "github.com/gin-gonic/gin"

type RouteConfig struct {
	App *gin.Engine
	Handler *Handler
}

func (rc *RouteConfig) Setup() {
	health := rc.App.Group("/health")

	health.GET("/live", rc.Handler.Live)
	health.GET("ready", rc.Handler.Ready)
}