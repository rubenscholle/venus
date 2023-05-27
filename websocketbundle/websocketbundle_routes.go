package websocketbundle

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitBundle(publicRoutes, protectedRoutes *gin.RouterGroup, ormDb *gorm.DB) {
	con := newWebsocketController(ormDb)

	protectedRoutes.GET("", con.WebsocketHandler)
}
