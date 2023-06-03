package websocketbundle

import (
	"github.com/gin-gonic/gin"
)

func (con *websocketController) WebsocketHandler(c *gin.Context) {
	manager.serveWs(c)
}
