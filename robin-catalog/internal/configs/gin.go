package configs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func NewGin(config *viper.Viper) *gin.Engine {
	var app = gin.Default()
	return app
}

func NewErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last()
			code := http.StatusInternalServerError

			var ginErr *gin.Error
			if errors.As(err.Err, &ginErr) {
				code = int(ginErr.Type)
			}

			ctx.JSON(code, gin.H{
				"errors": err.Error(),
			})
		}
	}
}
