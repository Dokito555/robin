package config

import (
	"github.com/Dokito555/robin-ums/internal/delivery/http"
	"github.com/Dokito555/robin-ums/internal/delivery/http/middleware"
	"github.com/Dokito555/robin-ums/internal/delivery/http/route"
	"github.com/Dokito555/robin-ums/internal/repository"
	"github.com/Dokito555/robin-ums/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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
	userRepository := repository.NewUserRepository(config.Log, config.DB)

	// setup services
	tokenService := services.NewTokenService(config.Log, config.Config)
	userService := services.NewUserService(config.DB, config.Log, config.Validate, userRepository, tokenService)

	// setup controllers
	healthController := http.NewHealthController(config.Log)

	// setup middleware
	middleware := middleware.NewAuth(userService, tokenService)

	// route config
	routeConfig := route.RouteConfig{
		App: config.App,
		HealthController: healthController,
		AuthMiddleware: middleware,
	}
	routeConfig.Setup()
}
