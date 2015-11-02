package wshub

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var (
	WRITE_WAIT  = 10 * time.Second
	PONG_WAIT   = 60 * time.Second
	PING_PERIOD = (PONG_WAIT * 9) / 10

	CONNECTION_MESSAGES_BUFFER_SIZE = 8

	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type Connection struct {
	socket  *websocket.Conn
	send    chan []byte
	receive chan []byte
	hub     *Hub
}

func (conn *Connection) readPump() {
	defer func() {
		if conn.hub != nil {
			conn.hub.UnregisterConnection(conn)
		}
		conn.socket.Close()
		close(conn.receive)
	}()

	conn.socket.SetReadDeadline(time.Now().Add(PONG_WAIT))
	conn.socket.SetPongHandler(func(string) error {
		conn.socket.SetReadDeadline(time.Now().Add(PONG_WAIT))
		return nil
	})

	for {
		if _, message, err := conn.socket.ReadMessage(); err == nil {
			conn.receive <- message
		} else {
			break
		}
	}
}

func (conn *Connection) write(messageType int, data []byte) error {
	conn.socket.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
	return conn.socket.WriteMessage(messageType, data)
}

func (conn *Connection) writePump() {
	ticker := time.NewTicker(PING_PERIOD)

	defer func() {
		conn.socket.Close()
		close(conn.send)
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-conn.send:
			if !ok {
				return
			}
			if err := conn.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := conn.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// Send message to connection
func (conn *Connection) SendMessage(message []byte) {
	conn.send <- message
}

// Incoming messages from connection
func (conn *Connection) Messages() <-chan []byte {
	return conn.receive
}

// Create new connection and start listen incoming messages
func NewConnection(rw http.ResponseWriter, request *http.Request) (*Connection, error) {
	socket, err := Upgrader.Upgrade(rw, request, nil)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		send:    make(chan []byte, CONNECTION_MESSAGES_BUFFER_SIZE),
		receive: make(chan []byte),
		socket:  socket,
	}

	go conn.writePump()
	go conn.readPump()

	return conn, nil
}
