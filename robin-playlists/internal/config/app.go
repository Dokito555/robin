package config

import (
	"github.com/Dokito555/robin-playlists/internal/delivery/grpc"
	"github.com/Dokito555/robin-playlists/internal/delivery/http"
	"github.com/Dokito555/robin-playlists/internal/delivery/http/middleware"
	"github.com/Dokito555/robin-playlists/internal/delivery/http/route"
	"github.com/Dokito555/robin-playlists/internal/repository"
	"github.com/Dokito555/robin-playlists/internal/services"
	"github.com/Dokito555/robin-playlists/pkg/ums"
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
	playlistRepo := repository.NewPlaylistRepository(config.Log, config.DB)

	// setup pkg
	ums := ums.NewUmsPkg(config.Log, config.Config)

	// setup services
	playlistService := services.NewPlaylistService(config.DB, config.Log, config.Validate, playlistRepo)
	tokenService := services.NewTokenService(config.Log, config.Config)

	// setup controllers
	healthController := http.NewHealthController(config.Log)
	playlistController := http.NewPlaylistController(config.Log, playlistService)

	// setup middleware
	middleware := middleware.NewAuth(ums, tokenService)

	// route config
	routeConfig := route.RouteConfig{
		App: config.App,
		HealthController: healthController,
		PlaylistController: playlistController,
		AuthMiddleware: middleware,
	}

	grpcConfig := grpc.GrpcConfig{
		Log: config.Log,
		Viper: config.Config,		
	}

	routeConfig.Setup()
	go grpcConfig.Setup()
}
