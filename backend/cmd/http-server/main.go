package main

import (
	"flag"
	"online-indicator/cmd/http-server/transport"
)

func main() {
	portUsage := "port to listen on"
	defaultPort := 8080
	port := flag.Int("port", defaultPort, portUsage)
	flag.IntVar(port, "p", defaultPort, portUsage)
	flag.Parse()
	server := transport.NewServer()
	err := server.Run(*port)
	if err != nil {
		panic(err)
	}
}
