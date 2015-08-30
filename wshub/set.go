package wshub

type Set struct {
	connections map[*Connection]bool
}

func (set *Set) Add(conn *Connection) {
	set.connections[conn] = true
}

func (set *Set) Remove(conn *Connection) {
	delete(set.connections, conn)
}

func (set *Set) Send(message *Message) {
	for conn := range set.connections {
		conn.send <- message.Text
	}
}

func NewSet() *Set {
	return &Set{
		connections: map[*Connection]bool{},
	}
}
