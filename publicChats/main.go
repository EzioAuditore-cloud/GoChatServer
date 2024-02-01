package main

import "gochatserver/server"

func main() {
	server := server.NewServer("127.0.0.1", 8080)
	server.Start()
}
