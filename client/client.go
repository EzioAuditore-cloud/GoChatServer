package client

import (
	gproto "gochatserver/model"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

func (c *Client) read(manager *ClientManager) {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		msg, _ := proto.Marshal(&gproto.ClientMessage{Sender: c.id, Content: string(message)})
		manager.broadcast <- msg
	}
}

func (c *Client) write(manager *ClientManager) {
	defer func() {
		c.socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
