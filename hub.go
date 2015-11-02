package wshub

import (
	"net/http"
)

var (
	HUB_MESSAGES_BUFFER_SIZE = 16
)

type Message struct {
	Text []byte
}

type Hub struct {
	connections *Set
	send        chan *Message
	register    chan *Connection
	unregister  chan *Connection
}

// Add conenction to hub
func (hub *Hub) RegisterConnection(conn *Connection) {
	conn.hub = hub
	hub.register <- conn
}

// Remove conenction from hub
func (hub *Hub) UnregisterConnection(conn *Connection) {
	hub.unregister <- conn
}

// Send message to all clients in hub
func (hub *Hub) SendMessage(message []byte) {
	hub.send <- &Message{
		Text: message,
	}
}

// Start handling connections and messages
func (hub *Hub) Run() {
	for {
		select {
		case conn := <-hub.register:
			hub.connections.Add(conn)
		case conn := <-hub.unregister:
			hub.connections.Remove(conn)
		case message := <-hub.send:
			hub.connections.Send(message)
		}
	}
}

// Create new conenction and add it to hub
func (hub *Hub) NewConnection(rw http.ResponseWriter, request *http.Request) (*Connection, error) {
	if conn, err := NewConnection(rw, request); err != nil {
		return nil, err
	} else {
		hub.RegisterConnection(conn)
		return conn, nil
	}

}

func NewHub() *Hub {
	return &Hub{
		send:        make(chan *Message, HUB_MESSAGES_BUFFER_SIZE),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: NewSet(),
	}
}
