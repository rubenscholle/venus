package websocketbundle

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Add frontend related events here
const (
	wsEventSendTextMessage = 1
	wsEventSendPongMessage = 10
)

type wsClient struct {
	connection *websocket.Conn
	manager    *wsManager

	// avoid multiple send processes at once
	sends chan []byte
}

type wsClientsMap map[*wsClient]bool

// wsManager holds reference to all clients connected via websocket and manages messages
type wsManager struct {
	clients wsClientsMap

	// ToDo use channels instead?
	sync.RWMutex

	// handlers are functions that are used to handle Events
	messageEvents map[uint]wsMessageEventHandler
}

type wsMessageEvent struct {
	MessageType uint   `json:"message_type"`
	Content     string `json:"content"`
}

type wsMessageEventHandler func(message wsMessageEvent, c *wsClient) error
