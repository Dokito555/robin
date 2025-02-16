package config

import (
	"github.com/Dokito555/robin-albums/internal/delivery/grpc"
	"github.com/Dokito555/robin-albums/internal/delivery/http"
	"github.com/Dokito555/robin-albums/internal/delivery/http/middleware"
	"github.com/Dokito555/robin-albums/internal/delivery/http/route"
	"github.com/Dokito555/robin-albums/internal/repository"
	"github.com/Dokito555/robin-albums/internal/services"
	"github.com/Dokito555/robin-albums/pkg/ums"
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
	albumRepository := repository.NewAlbumRepository(config.Log, config.DB)

	// setup services
	tokenService := services.NewTokenService(config.Log, config.Config)
	albumService := services.NewAlbumService(config.DB, config.Log, config.Validate, albumRepository)

	// setup pkg
	ums := ums.NewUmsPkg(config.Log, config.Config)

	// setup controllers
	healthController := http.NewHealthController(config.Log)
	albumController := http.NewAlbumService(config.Log, albumService)

	// setup middleware
	authMiddleware := middleware.NewAuth(ums, tokenService)

	// route config
	routeConfig := route.RouteConfig{
		App:              config.App,
		HealthController: healthController,
		AlbumController:  albumController,
		AuthMiddleware:   authMiddleware,
	}
	
	// setup grpc
	grpcConfig := grpc.GrpcConfig{
		Log: config.Log,
		Viper: config.Config,
	}

	routeConfig.Setup()
	go grpcConfig.Setup()
}
