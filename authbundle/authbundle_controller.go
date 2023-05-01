package authbundle

import (
	core "github.com/rubenscholle/venus/corebundle"
	"gorm.io/gorm"
)

type authController struct {
	core.Controller
	OrmDb gorm.DB
}

func newAuthController(ormDb *gorm.DB) *authController {
	con := &authController{
		OrmDb: *ormDb,
	}

	ormDb.AutoMigrate(&AuthUser{})

	return con
}
