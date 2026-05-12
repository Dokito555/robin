package config

import (
	"github.com/Dokito555/robin/internal/healthcheck"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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
}

func Bootstrap(config *BootstrapConfig) {
	// health
	healthcheck.Register(config.App, config.DB)
}
