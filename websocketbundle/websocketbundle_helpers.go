package websocketbundle

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	m.messageEvents[wsEventSendMessage] = func(message wsMessageEvent, c *wsClient) error {
		// ToDo
		log.Println(message)
		return nil
	}
}

// routeWsMessageEvent is used to make sure the correct event goes into the correct handler
func (m *wsManager) routeWsMessageEvent(message wsMessageEvent, c *wsClient) error {
	// Check if Handler is present in Map
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

	// Check if Client exists, then delete it
	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

func (c *wsClient) readMessage() {
	defer c.manager.removeClient(c)

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			// to not log disconnects
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
	defer c.manager.removeClient(c)

	for {

		// add other cases later
		select {
		case rawMessage, ok := <-c.sends:
			if !ok {
				// Manager has closed connection channel
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println(err)
				}

				return
			}

			message, err := json.Marshal(rawMessage)
			if err != nil {
				log.Println(err)
				return // closes the connection, should we really
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
			}
		}
	}
}
