package main

import (
	"flag"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/inmem"
	"github.com/Kshitij09/online-indicator/transport"
	"github.com/jonboulle/clockwork"
	"os"
	"strconv"
	"time"
)

var DefaultConfig = domain.Config{
	OnlineThreshold: 10 * time.Second,
	ServerPort:      8080,
}

func intEnvOrDefault(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		value, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return value
	}
	return fallback
}

func main() {
	cfg := DefaultConfig
	flag.IntVar(&cfg.ServerPort, "p", intEnvOrDefault("PORT", 8080), "port to listen on")
	envOnlineThreshold := time.Duration(intEnvOrDefault("ONLINE_THRESHOLD_MILLIS", 10_000)) * time.Millisecond
	flag.DurationVar(&cfg.OnlineThreshold, "online-threshold", envOnlineThreshold, "threshold for determining if a user is online")
	flag.Parse()

	tokenGen := domain.NewUUIDTokenGenerator()
	sessionGen := domain.NewUUIDSessionGenerator()
	realClock := clockwork.NewRealClock()
	idGen := domain.NewSeqIdGenerator()
	storage := inmem.NewStorage(tokenGen, sessionGen, realClock, idGen)
	server := transport.NewServer(storage, cfg, realClock)
	err := server.Run(cfg.ServerPort)
	if err != nil {
		panic(err)
	}
}
