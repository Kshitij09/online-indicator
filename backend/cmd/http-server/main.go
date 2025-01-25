package main

import (
	"flag"
	"github.com/Kshitij09/online-indicator/cmd/http-server/transport"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/inmem"
)

func main() {
	portUsage := "port to listen on"
	defaultPort := 8080
	port := flag.Int("port", defaultPort, portUsage)
	flag.IntVar(port, "p", defaultPort, portUsage)
	flag.Parse()
	tokenGen := domain.NewUUIDTokenGenerator()
	storage := inmem.NewStorage(tokenGen)
	server := transport.NewServer(storage)
	err := server.Run(*port)
	if err != nil {
		panic(err)
	}
}
