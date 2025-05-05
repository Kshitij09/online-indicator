package transport

import (
	"context"
	"fmt"
	"github.com/Kshitij09/online-indicator/di"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/redisstore"
	"github.com/Kshitij09/online-indicator/transport/middlewares"
	"github.com/jonboulle/clockwork"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

type Server struct {
	domain.Storage
	config      domain.Config
	clock       clockwork.Clock
	lastSeenDao domain.LastSeenDao
	db          di.DatabaseContainer
	svcs        di.ServiceContainer
	handlers    di.HandlerContainer
}

func NewServer(storage domain.Storage, config domain.Config, clock clockwork.Clock, client *redis.Client, db di.DatabaseContainer, svcs di.ServiceContainer, handlers di.HandlerContainer) *Server {
	return &Server{
		Storage:     storage,
		config:      config,
		clock:       clock,
		lastSeenDao: redisstore.LastSeenDao(client, context.Background(), config.OnlineThreshold),
		db:          db,
		svcs:        svcs,
		handlers:    handlers,
	}
}
func (s *Server) Run(port int) error {
	listAddr := fmt.Sprintf(":%d", port)
	router := http.NewServeMux()
	logger := middlewares.HttpLogger
	router.HandleFunc("GET /health", NewHttpHandler(func(w http.ResponseWriter, r *http.Request) error {
		_, err := w.Write([]byte("OK"))
		return err
	}, logger))
	router.HandleFunc("POST /register", NewHttpHandler(s.handlers.Register, logger))
	router.HandleFunc("POST /login", NewHttpHandler(s.handlers.Login, logger))
	router.HandleFunc(fmt.Sprintf("POST /ping/{%s}", PathId), NewHttpHandler(s.handlers.Ping, logger))
	router.HandleFunc(fmt.Sprintf("GET /status/{%s}", PathId), NewHttpHandler(s.handlers.Status, logger))
	router.HandleFunc("POST /batch/status", NewHttpHandler(s.handlers.BatchStatus, logger))
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
