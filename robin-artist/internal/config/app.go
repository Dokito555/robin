package config

import (
	"github.com/Dokito555/robin-artist/internal/delivery/http"
	"github.com/Dokito555/robin-artist/internal/delivery/http/middleware"
	"github.com/Dokito555/robin-artist/internal/delivery/http/route"
	"github.com/Dokito555/robin-artist/internal/repository"
	"github.com/Dokito555/robin-artist/internal/services"
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
	artistRepo := repository.NewAristRepository(config.Log, config.DB)

	// setup services
	tokenService := services.NewTokenService(config.Log, config.Config)
	artistService := services.NewArtistService(config.DB, config.Log, config.Validate, artistRepo, tokenService)

	// setup controllers
	healthController := http.NewHealthController(config.Log)
	artistController := http.NewArtistController(config.Log, artistService)

	// setup middleware
	middleeware := middleware.NewAuth(artistService, tokenService)

	// route config
	routeConfig := route.RouteConfig{
		App: config.App,
		HealthController: healthController,
		ArtistController: artistController,
		AuthMiddleware: middleeware,
	}
	
	routeConfig.Setup()
}
