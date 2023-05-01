package authbundle

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitBundle(publicRoutes, protectedRoutes *gin.RouterGroup, ormDb *gorm.DB) {
	con := newAuthController(ormDb)

	publicRoutes.POST("/login", con.LoginHandler)
}
