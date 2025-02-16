package config

import (
	"github.com/Dokito555/robin-playlists/internal/delivery/http"
	"github.com/Dokito555/robin-playlists/internal/delivery/http/route"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *gin.Engine
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup repo

	// setup services

	// setup controllers
	healthController := http.NewHealthController(config.Log)

	// setup middleware

	// route config
	routeConfig := route.RouteConfig{
		App: config.App,
		HealthController: healthController,
	}
	routeConfig.Setup()
}
