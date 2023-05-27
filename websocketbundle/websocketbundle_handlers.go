package websocketbundle

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (con *websocketController) WebsocketHandler(c *gin.Context) {
	wsCon, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}
	defer wsCon.Close()

	messageType, _, err := wsCon.ReadMessage()
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}

	response := []byte("Hello, client!")
	err = wsCon.WriteMessage(messageType, response)
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}
}
