package config

import (
	"github.com/Dokito555/robin-notification/internal/delivery/http"
	"github.com/Dokito555/robin-notification/internal/delivery/http/route"
	"github.com/Dokito555/robin-notification/internal/repository"
	"github.com/Dokito555/robin-notification/internal/services"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB            *gorm.DB
	App           *gin.Engine
	Log           *logrus.Logger
	Validate      *validator.Validate
	Config        *viper.Viper
	KafkaConsumer sarama.Consumer
}

func Bootstrap(config *BootstrapConfig) {
	// setup repo
	emailRepository := repository.NewEmailRepository(config.Log, config.DB)
	templateRepository := repository.NewTemplateRepository(config.Log, config.DB)

	// setup services
	messagingService := services.NewMessagingService(config.Log, config.KafkaConsumer)
	notificationService := services.NewNotificationService(config.Log, emailRepository, templateRepository, messagingService)

	// setup controllers
	healthController := http.NewHealthController(config.Log)
	notificationController := http.NewNotificationController(notificationService)

	// setup middleware

	// route config
	routeConfig := route.RouteConfig{
		App:                    config.App,
		HealthController:       healthController,
		NotificationController: notificationController,
	}
	routeConfig.Setup()
}
