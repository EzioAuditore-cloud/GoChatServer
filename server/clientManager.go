package server

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Client struct {
	ID         int32
	UUID       string
	Addr       string
	Name       string
	Socket     *websocket.Conn
	Conn       net.Conn
	SendBytes  chan []byte
	SendString chan string
	ChatMap    map[*Client][]string
	Srv        *Server
	State      int32 //0-online 1-offline
}

func NewClient(conn interface{}, server *Server) *Client {
	clientId++
	client := &Client{
		ID:         int32(clientId),
		UUID:       uuid.NewV4().String(),
		SendBytes:  make(chan []byte),
		SendString: make(chan string),
		Srv:        server,
	}
	switch val := conn.(type) {
	case *websocket.Conn:
		client.Socket = val
		client.Addr = val.RemoteAddr().String()
	case net.Conn:
		client.Conn = val
		client.Addr = val.RemoteAddr().String()
	default:
		fmt.Println("NewClient unknow conn")
		return nil
	}
	server.mapLock.Lock()
	server.ClientMap[client.Addr] = client
	server.mapLock.Unlock()
	go client.ListenSend()

	return client
}

func (c *Client) ListenSend() {
	for {
		stringMsg := <-c.SendString
		c.Conn.Write([]byte(stringMsg + "\n"))
	}
}

func Login(conn net.Conn, srv *Server) *Client {
	var client *Client
	srv.mapLock.Lock()
	if v, ok := srv.ClientMap[conn.RemoteAddr().String()]; !ok {
		client = NewClient(conn, srv)
	} else {
		client = v
	}
	srv.mapLock.Unlock()
	srv.BroadCast(client, "已上线")
	return client
}

func (c *Client) Logout() {
	srv := c.Srv
	srv.BroadCast(c, "已下线")
	c.State = 1
}

func (c *Client) DoMessage(conn net.Conn, buf []byte, ctx context.Context, isLive chan bool) {
	srv := c.Srv
	for {
		select {
		case <-ctx.Done():
			c.Logout()
			return
		default:
			n, err := conn.Read(buf)
			if n == 0 {
				c.Logout()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read to buf err:", err)
				return
			}
			msg := string(buf[:n-1])
			srv.BroadCast(c, msg)
			isLive <- true
		}
	}
}

// func (c *Client) Read(manager *ClientManager) {
// 	defer func() {
// 		manager.Unregister <- c
// 		c.Socket.Close()
// 	}()
// 	for {
// 		_, message, err := c.Socket.ReadMessage()
// 		if err != nil {
// 			manager.Unregister <- c
// 			c.Socket.Close()
// 			break
// 		}
// 		msg, _ := proto.Marshal(&gproto.ClientMessage{Sender: c.Id, Content: string(message)})
// 		manager.Broadcast <- msg
// 	}
// }

// func (c *Client) Write(manager *ClientManager) {
// 	defer func() {
// 		c.Socket.Close()
// 	}()
// 	for {
// 		select {
// 		case message, ok := <-c.SendBytes:
// 			if !ok {
// 				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}
// 			c.Socket.WriteMessage(websocket.TextMessage, message)
// 		}
// 	}
// }
