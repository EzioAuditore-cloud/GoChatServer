package client

import (
	gproto "gochatserver/model"

	"google.golang.org/protobuf/proto"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (manager *ClientManager) send(msg []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- msg
		}
	}
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			msgBytes, _ := proto.Marshal(&gproto.ClientMessage{Content: "/A new socket has connected."})
			manager.send(msgBytes, conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				// manager.clients[conn] = false
				msgBytes, _ := proto.Marshal(&gproto.ClientMessage{Content: "/A new socket has disconnected."})
				manager.send(msgBytes, conn)
			}
		case msg := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- msg:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}
