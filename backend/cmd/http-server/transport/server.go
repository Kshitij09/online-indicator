package transport

import (
	"fmt"
	"github.com/Kshitij09/online-indicator/cmd/http-server/transport/handlers"
	"github.com/Kshitij09/online-indicator/cmd/http-server/transport/middlewares"
	"log"
	"net/http"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}
func (s *Server) Run(port int) error {
	listAddr := fmt.Sprintf(":%d", port)
	router := http.NewServeMux()
	logger := middlewares.HttpLogger
	router.HandleFunc("GET /health", handlers.NewHttpHandler(health, logger))
	server := &http.Server{
		Addr:    listAddr,
		Handler: router,
	}
	log.Println("online indicator server started listening on " + listAddr)
	return server.ListenAndServe()
}

func health(w http.ResponseWriter, _ *http.Request) error {
	_, err := w.Write([]byte("OK"))
	return err
}
