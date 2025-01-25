package transport

import (
	"fmt"
	"github.com/Kshitij09/online-indicator/cmd/http-server/transport/middlewares"
	"github.com/Kshitij09/online-indicator/domain"
	"log"
	"net/http"
)

type Server struct {
	domain.Storage
}

func NewServer(storage domain.Storage) *Server {
	return &Server{
		Storage: storage,
	}
}
func (s *Server) Run(port int) error {
	listAddr := fmt.Sprintf(":%d", port)
	router := http.NewServeMux()
	logger := middlewares.HttpLogger
	router.HandleFunc("GET /health", NewHttpHandler(health, logger))
	register := RegisterHandler(s.Storage)
	router.HandleFunc("POST /register", NewHttpHandler(register, logger))
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
