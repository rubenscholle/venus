package websocketbundle

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	core "github.com/rubenscholle/venus/corebundle"
)

const (
	pongWait = 10 * time.Second

	pingInterval = (pongWait * 9) / 10
)

// newWsClient initialites a wsClient object representing a frontend client connected via websocket
func newWsClient(con *websocket.Conn, manager *wsManager) *wsClient {
	return &wsClient{
		connection: con,
		manager:    manager,
		sends:      make(chan []byte),
	}
}

// NewWsManager initializes wsManager object that holds a map with a reference to each client connected
func newWsManager() *wsManager {
	m := &wsManager{
		clients:       make(wsClientsMap),
		messageEvents: make(map[uint]wsMessageEventHandler),
	}
	m.setupEventHandlers()

	return m
}

// serveWs is called by the websocketHandler and handles websocket requests via http
func (m *wsManager) serveWs(c *gin.Context) {
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}
	//defer wsConn.Close()

	client := newWsClient(wsConn, m)

	m.addClient(client)

	go client.readMessage()
	go client.writeMessages()
}

// setupEventHandlers configures and adds all handlers
func (m *wsManager) setupEventHandlers() {
	m.messageEvents[wsEventSendTextMessage] = func(message wsMessageEvent, c *wsClient) error {
		// ToDo
		log.Println("text message received via websocket")
		return nil
	}
}

// routeWsMessageEvent is used to make sure the correct event goes into the correct handler
func (m *wsManager) routeWsMessageEvent(message wsMessageEvent, c *wsClient) error {
	// check if Handler is present in Map
	if handler, ok := m.messageEvents[message.MessageType]; ok {
		if err := handler(message, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("event not supported")
	}
}

func (m *wsManager) addClient(client *wsClient) {
	// ToDo use channels instead?
	// only one process can read a map at a time
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *wsManager) removeClient(client *wsClient) {
	// ToDo use channels instead?
	// only one process can read a map at a time
	m.Lock()
	defer m.Unlock()

	// check if Client exists, then delete it
	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

func (c *wsClient) readMessage() {
	defer c.manager.removeClient(c)

	// limit long websocket messages
	c.connection.SetReadLimit(512)

	// set initial ping timer
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			// disconnects should not be logged
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}

		var messageEvent wsMessageEvent
		if err := json.Unmarshal(payload, &messageEvent); err != nil {
			log.Println(err)
			break
		}

		if err := c.manager.routeWsMessageEvent(messageEvent, c); err != nil {
			log.Println(err)
		}
	}
}

func (c *wsClient) writeMessages() {
	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()

		c.manager.removeClient(c)
	}()

	for {
		select {
		case rawMessage, ok := <-c.sends:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println(err)
				}

				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, rawMessage); err != nil {
				log.Println(err)
			}

		case <-ticker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// pongHandler resets wait time before closing the websocket every time the client reponds with a pong
func (c *wsClient) pongHandler(pongMsg string) error {
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

// deprecated
// checkOrigin check if websocket connection origin is valid
func checkOrigin(r *http.Request) bool {

	// Grab the request origin
	origin := r.Header.Get("Origin")

	switch origin {
	case fmt.Sprintf("ws://localhost:%d", core.Config.Server.Port), fmt.Sprintf("wss://localhost:%d", core.Config.Server.Port):
		return true
	default:
		return false
	}
}

func BroadcastEventMessage(message string) {
	currentClients := manager.clients

	for client, _ := range currentClients {
		client.sends <- []byte(message)
	}
}
