package websocketbundle

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// newWsClient initialites a wsClient object representing a frontend client connected via websocket
func newWsClient(con *websocket.Conn, manager *wsManager) *wsClient {
	return &wsClient{
		connection: con,
		manager:    manager,
	}
}

// NewWsManager initializes wsManager object that holds a map with a reference to each client connected
func newWsManager() *wsManager {
	return &wsManager{
		clients: make(wsClientsMap),
	}
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

	messageType, _, err := wsConn.ReadMessage()
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}

	response := []byte("Hello, client!")
	err = wsConn.WriteMessage(messageType, response)
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}

	log.Println(len(m.clients))

	// Todo
	//go client.readMessages()
	//go client.writeMessages()
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
