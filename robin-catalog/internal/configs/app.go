package config

import (
	"github.com/Dokito555/robin/robin-catalog/internal/delivery/grpc"
	"github.com/Dokito555/robin/robin-catalog/internal/delivery/http"
	"github.com/Dokito555/robin/robin-catalog/internal/delivery/http/middleware"
	"github.com/Dokito555/robin/robin-catalog/internal/delivery/http/route"
	"github.com/Dokito555/robin/robin-catalog/internal/repository"
	"github.com/Dokito555/robin/robin-catalog/internal/services"
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
	albumRepo := repository.NewAlbumRepository(config.Log, config.DB)
	// songRepository := repository.NewSongRepository(config.Log, config.DB)

	// setup services
	tokenService := services.NewTokenService(config.Log, config.Config)
	artistService := services.NewArtistService(config.DB, config.Log, config.Validate, artistRepo, tokenService)
	albumService := services.NewAlbumService(config.DB, config.Log, config.Validate, albumRepo)
	// artistPkg := artist.NewArtist(config.Log, config.Config)
	// songService := services.NewSongService(config.Log, config.DB, config.Validate, config.Config, songRepository)

	// setup controllers
	healthController := http.NewHealthController(config.Log)
	artistController := http.NewArtistController(config.Log, artistService)
	albumController := http.NewAlbumService(config.Log, albumService)
	// songController := http.NewSongController(config.Log, songService)
	tokenValidationController := grpc.NewTokenValidationController(tokenService, config.Log)

	// setup middleware
	middleeware := middleware.NewAuth(artistService, tokenService)

	// route config
	routeConfig := route.RouteConfig{
		App:              config.App,
		HealthController: healthController,
		ArtistController: artistController,
		AlbumController:  albumController,
		// SongController:   songController,
		AuthMiddleware:   middleeware,
	}

	grpcConfig := grpc.GrpcConfig{
		Log:                       config.Log,
		Viper:                     config.Config,
		TokenValidationController: tokenValidationController,
	}

	routeConfig.Setup()
	go grpcConfig.Setup()
}
