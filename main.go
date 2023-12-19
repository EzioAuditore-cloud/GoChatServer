package main

import (
	"fmt"
	"net/http"

	gproto "gochatserver/model"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

var manager = ClientManager{
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
	broadcast:  make(chan []byte),
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

func (c *Client) read() {
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
		fmt.Println("Client id:", c.id)
		msg, _ := proto.Marshal(&gproto.ClientMessage{Sender: c.id, Content: string(message)})
		manager.broadcast <- msg
	}
}

func (c *Client) write() {
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

func main() {

	fmt.Println("Starting application...")
	go manager.start()
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":12345", nil)

}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		fmt.Println("no found")
		return
	}

	client := &Client{id: uuid.NewV4().String(), socket: conn, send: make(chan []byte)}
	manager.register <- client
	go client.read()
	go client.write()
}
