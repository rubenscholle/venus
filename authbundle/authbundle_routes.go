package authbundle

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitBundle(rg *gin.RouterGroup, ormDb *gorm.DB) {
	con := newAuthController(ormDb)

	rg.POST("/login", con.LoginHandler)
}
