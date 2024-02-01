package server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port      int
	ClientMap map[string]*Client //客户端map key-Addr value-Client
	mapLock   sync.RWMutex

	//广播消息channel
	BroadcastChannel chan string
}

var clientId = 0

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:               ip,
		Port:             port,
		ClientMap:        make(map[string]*Client),
		BroadcastChannel: make(chan string),
	}

	return server
}

func (srv *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", srv.Ip, srv.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	fmt.Printf("服务启动: %s:%d\n", srv.Ip, srv.Port)
	defer listener.Close()
	go srv.ListenMessage()
	for {
		conn, err := listener.Accept() //建立连接
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}
		go srv.Handler(conn) //启动goroutine处理连接
	}

}

//监听广播channel
func (srv *Server) ListenMessage() {
	for {
		msg := <-srv.BroadcastChannel
		srv.mapLock.Lock()
		for _, cli := range srv.ClientMap {
			if cli.State == 0 {
				cli.SendString <- msg
			}
		}
		srv.mapLock.Unlock()
	}
}

func (srv *Server) BroadCast(client *Client, msg string) {
	sendMsg := "[" + client.Addr + "]" + client.UUID + ":" + msg
	srv.BroadcastChannel <- sendMsg
}

func (srv *Server) Login(conn net.Conn) *Client {
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

func (srv *Server) Handler(conn net.Conn) {
	client := srv.Login(conn)
	fmt.Printf("%v连接成功\n", client.ID)
	defer func() {
		e := recover()
		if e != nil {
			fmt.Println("panic: ", e)
		}
		conn.Close()
		fmt.Printf("%v连接断开\n", client.ID)
	}()

	buf := make([]byte, 1024)
	isLive := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	go func(context.Context) {
		client.DoMessage(conn, buf, ctx, isLive)
	}(ctx)

	for {
		select {
		case <-isLive:
		case <-time.After(time.Second * 10):
			cancel()
			return
		}
	}
}
