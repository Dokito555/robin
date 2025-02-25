package config

import (
	"github.com/Dokito555/robin-songs/internal/delivery/grpc"
	"github.com/Dokito555/robin-songs/internal/delivery/http"
	"github.com/Dokito555/robin-songs/internal/delivery/http/middleware"
	"github.com/Dokito555/robin-songs/internal/delivery/http/route"
	"github.com/Dokito555/robin-songs/internal/repository"
	"github.com/Dokito555/robin-songs/internal/services"
	"github.com/Dokito555/robin-songs/pkg/artist"
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
	S3Client *S3Client
}

func Bootstrap(config *BootstrapConfig) {
	// setup repo
	songRepository := repository.NewSongRepository(config.Log, config.DB)

	// setup pkg
	artistPkg := artist.NewArtist(config.Log, config.Config)

	// setup services
	tokenService := services.NewTokenService(config.Log, config.Config)
	songService := services.NewSongService(config.Log, config.DB, config.Validate, config.Config, songRepository)

	// setup controllers
	healthController := http.NewHealthController(config.Log)
	songController := http.NewSongController(config.Log, songService)

	// setup middleware
	authMiddleware := middleware.NewAuth(artistPkg, tokenService)

	// route config
	routeConfig := route.RouteConfig{
		App: config.App,
		HealthController: healthController,
		SongController: songController,
		AuthMiddleware: authMiddleware,
	}

	// grpc config
	grpcConfig := grpc.GrpcConfig{
		Log: config.Log,
		Viper: config.Config,
	}

	go grpcConfig.RunGrpc()

	routeConfig.Setup()
}
