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
	config domain.Config
}

func NewServer(storage domain.Storage, config domain.Config) *Server {
	return &Server{
		Storage: storage,
		config:  config,
	}
}
func (s *Server) Run(port int) error {
	listAddr := fmt.Sprintf(":%d", port)
	router := http.NewServeMux()
	logger := middlewares.HttpLogger
	router.HandleFunc("GET /health", NewHttpHandler(health, logger))
	register := RegisterHandler(s.Storage)
	router.HandleFunc("POST /register", NewHttpHandler(register, logger))
	login := LoginHandler(s.Storage)
	router.HandleFunc("POST /login", NewHttpHandler(login, logger))
	ping := PingHandler(s.Storage, s.config)
	router.HandleFunc("POST /ping", NewHttpHandler(ping, logger))
	status := StatusHandler(s.Storage, s.config)
	router.HandleFunc(fmt.Sprintf("GET /status/{%s}", PathId), NewHttpHandler(status, logger))
	batchStatus := BatchStatusHandler(s.Storage, s.config)
	router.HandleFunc("POST /batch/status", NewHttpHandler(batchStatus, logger))
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
