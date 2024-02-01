package app

// import (
// 	"context"
// 	"fmt"
// 	"gochatserver/client"
// 	"io"
// 	"net"
// 	"sync"
// 	"time"
// )

// type Server struct {
// 	Ip        string
// 	Port      int
// 	ClientMap map[string]*client.Client //客户端map key-id value-Client
// 	mapLock   sync.RWMutex

// 	//广播消息channel
// 	BroadcastChannel chan string
// }

// var clientId = 0

// func NewServer(ip string, port int) *Server {
// 	server := &Server{
// 		Ip:               ip,
// 		Port:             port,
// 		ClientMap:        make(map[string]*client.Client),
// 		BroadcastChannel: make(chan string),
// 	}

// 	return server
// }

// func (srv *Server) Start() {
// 	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", srv.Ip, srv.Port))
// 	if err != nil {
// 		fmt.Println("net.Listen err:", err)
// 		return
// 	}
// 	fmt.Printf("服务启动: %s:%d\n", srv.Ip, srv.Port)
// 	defer listener.Close()
// 	go srv.ListenMessage()
// 	for {
// 		conn, err := listener.Accept() //建立连接
// 		if err != nil {
// 			fmt.Println("listener.Accept err:", err)
// 			continue
// 		}
// 		go srv.Handler(conn) //启动goroutine处理连接
// 	}

// }

// //监听广播channel
// func (srv *Server) ListenMessage() {
// 	for {
// 		msg := <-srv.BroadcastChannel
// 		srv.mapLock.Lock()
// 		for _, cli := range srv.ClientMap {
// 			cli.SendString <- msg
// 		}
// 		srv.mapLock.Unlock()
// 	}
// }

// func (srv *Server) BroadCast(client *client.Client, msg string) {
// 	sendMsg := "[" + client.Addr + "]" + client.Id + ":" + msg
// 	srv.BroadcastChannel <- sendMsg
// }

// func (srv *Server) Handler(conn net.Conn) {
// 	clientId++
// 	id := clientId
// 	fmt.Printf("%v连接成功\n", id)
// 	defer func() {
// 		e := recover()
// 		if e != nil {
// 			fmt.Println("panic: ", e)
// 		}
// 		conn.Close()
// 		fmt.Printf("%v连接断开\n", id)
// 	}()
// 	client := client.NewClient(conn)
// 	srv.mapLock.Lock()
// 	srv.ClientMap[client.Id] = client
// 	srv.mapLock.Unlock()
// 	srv.BroadCast(client, "已上线")
// 	buf := make([]byte, 1024)
// 	isLive := make(chan bool)
// 	ctx, cancel := context.WithCancel(context.Background())
// 	go func(context.Context) {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				srv.BroadCast(client, "已下线")
// 				delete(client.ChatMap, client)
// 				return
// 			default:
// 				n, err := conn.Read(buf)
// 				if n == 0 {
// 					srv.BroadCast(client, "已下线")
// 					delete(client.ChatMap, client)
// 					return
// 				}
// 				if err != nil && err != io.EOF {
// 					fmt.Println("conn.Read to buf err:", err)
// 					return
// 				}
// 				msg := string(buf[:n-1])
// 				srv.BroadCast(client, msg)
// 				isLive <- true
// 			}
// 		}
// 	}(ctx)

// 	for {
// 		select {
// 		case <-isLive:
// 		case <-time.After(time.Second * 10):
// 			cancel()
// 			return
// 		}
// 	}
// }
