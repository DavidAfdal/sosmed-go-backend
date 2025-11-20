package socket

import (
	"github.com/gorilla/websocket"
)

type Connection struct {
	ws   *websocket.Conn
	send chan []byte
}

func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		ws:   ws,
		send: make(chan []byte, 256),
	}
}

func (c *Connection) Send(msg []byte) {
	select {
	case c.send <- msg:
	default:
		// drop or close connection to avoid memory leak
	}
}

// func (c *Connection) writePump() {
// 	ticker := time.NewTicker(time.Second * 50)

// 	defer func() {
// 		ticker.Stop()
// 		c.ws.Close()
// 	}()

// 	for {
// 		select {
// 		// case msg, ok := <- c.send:

// 		// }
// 	}
// }
