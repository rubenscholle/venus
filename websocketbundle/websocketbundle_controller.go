package websocketbundle

import (
	"github.com/gorilla/websocket"
	core "github.com/rubenscholle/venus/corebundle"
	"gorm.io/gorm"
)

type websocketController struct {
	core.Controller
	OrmDb gorm.DB
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func newWebsocketController(ormDb *gorm.DB) *websocketController {
	con := &websocketController{
		OrmDb: *ormDb,
	}

	//ormDb.AutoMigrate(&{})

	return con
}
