package config

import (
	"github.com/Dokito555/robin-ums/constants"
	"github.com/Dokito555/robin-ums/internal/delivery/grpc"
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
	artistRepository := repository.NewAristRepository(config.Log, config.DB)

	// setup services
	tokenService := services.NewTokenService(config.Log, config.Config)
	userService := services.NewUserService(config.DB, config.Log, config.Validate, userRepository, tokenService)
	artistService := services.NewArtistService(config.DB, config.Log, config.Validate, artistRepository, tokenService)

	// setup controllers
	healthController := http.NewHealthController(config.Log)
	userController := http.NewUserController(config.Log, userService)
	artistController := http.NewArtistController(config.Log, artistService)
	tokenValidationController := grpc.NewTokenValidationController(tokenService, config.Log)

	// setup middleware
	userMiddleware := middleware.NewAuth(userService, tokenService, constants.ROLE_USER)
	adminMiddleware := middleware.NewAuth(userService, tokenService, constants.ROLE_ADMIN)
	artistMiddleware := middleware.NewAuth(userService, tokenService, constants.ROLE_ARTIST)

	// route config
	routeConfig := route.RouteConfig{
		App:              config.App,
		HealthController: healthController,
		UserController:   userController,
		ArtistController: artistController,
		UserMiddleware:   userMiddleware,
		AdminMiddleware:  adminMiddleware,
		ArtistMiddelware: artistMiddleware,
	}

	// grpc config
	grpcConfig := grpc.GrpcConfig{
		Log:                       config.Log,
		Viper:                     config.Config,
		TokenValidationController: tokenValidationController,
	}

	routeConfig.Setup()
	go grpcConfig.Setup()
}
