package client

// import (
// 	gproto "gochatserver/model"

// 	"google.golang.org/protobuf/proto"
// )

// type ClientManager struct {
// 	Clients    map[*Client]bool
// 	Broadcast  chan []byte
// 	Register   chan *Client
// 	Unregister chan *Client
// }

// func (manager *ClientManager) Send(msg []byte, ignore *Client) {
// 	for conn := range manager.Clients {
// 		if conn != ignore {
// 			conn.SendBytes <- msg
// 		}
// 	}
// }

// func (manager *ClientManager) Start() {
// 	for {
// 		select {
// 		case conn := <-manager.Register:
// 			manager.Clients[conn] = true
// 			msgBytes, _ := proto.Marshal(&gproto.ClientMessage{Content: "/A new socket has connected."})
// 			manager.Send(msgBytes, conn)
// 		case conn := <-manager.Unregister:
// 			if _, ok := manager.Clients[conn]; ok {
// 				close(conn.SendBytes)
// 				delete(manager.Clients, conn)
// 				// manager.clients[conn] = false
// 				msgBytes, _ := proto.Marshal(&gproto.ClientMessage{Content: "/A new socket has disconnected."})
// 				manager.Send(msgBytes, conn)
// 			}
// 		case msg := <-manager.Broadcast:
// 			for conn := range manager.Clients {
// 				select {
// 				case conn.SendBytes <- msg:
// 				default:
// 					close(conn.SendBytes)
// 					delete(manager.Clients, conn)
// 				}
// 			}
// 		}
// 	}
// }
