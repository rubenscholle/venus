package corebundle

import (
	"time"

	"gorm.io/gorm"
)

type Controller struct{}

type SystemConfiguration struct {
	Server   ServerConfiguration   `json:"server"`
	Database DatabaseConfiguration `json:"database"`
}

type ServerConfiguration struct {
	TokenTTL      uint   `json:"token_ttl"`
	JWTPrivateKey string `json:"jwt_private_key"`
}

type DatabaseConfiguration struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
