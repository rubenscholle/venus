package websocketbundle

import (
	"sync"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	connection *websocket.Conn
	manager    *wsManager
}

type wsClientsMap map[*wsClient]bool

// wsManager holds reference to all clients connected via websocket and manages messages
type wsManager struct {
	clients wsClientsMap

	// ToDo use channels instead?
	sync.RWMutex
}
