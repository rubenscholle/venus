package authbundle

import (
	core "github.com/rubenscholle/venus/corebundle"
	"gorm.io/gorm"
)

// ToDo make private
type AuthController struct {
	core.Controller
	OrmDb gorm.DB
}
