package healthcheck

import "gorm.io/gorm"
import "github.com/gin-gonic/gin"

func Register(app *gin.Engine, db *gorm.DB) {
	service := NewService(db)
	handler := NewHandler(service)

	routes := RouteConfig{
		App: app,
		Handler: handler,
	}

	routes.Setup()
}