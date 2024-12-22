package transport

import (
	"fmt"
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
	router.HandleFunc("GET /health", health)
	server := &http.Server{
		Addr:    listAddr,
		Handler: router,
	}
	log.Println("online indicator server started listening on " + listAddr)
	return server.ListenAndServe()
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}
