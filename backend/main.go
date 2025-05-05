package main

import (
	"context"
	"flag"
	"github.com/Kshitij09/online-indicator/di"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/service"
	"github.com/Kshitij09/online-indicator/inmem"
	"github.com/Kshitij09/online-indicator/redisstore"
	"github.com/Kshitij09/online-indicator/transport"
	"github.com/jonboulle/clockwork"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

var DefaultConfig = domain.Config{
	OnlineThreshold: 10 * time.Second,
	ServerPort:      8080,
}

func main() {
	cfg := parseConfig()
	apiKeyGen := domain.NewUUIDApiKeyGenerator()
	sessionGen := domain.NewUUIDSessionGenerator()
	realClock := clockwork.NewRealClock()
	idGen := domain.NewSeqIdGenerator()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})
	db := di.DatabaseContainer{
		Auth:     inmem.NewAuthDao(apiKeyGen, idGen),
		Session:  inmem.NewSessionDao(sessionGen, realClock),
		Profile:  inmem.NewProfileDao(),
		LastSeen: redisstore.LastSeenDao(redisClient, context.Background(), cfg.OnlineThreshold),
	}
	svcs := di.ServiceContainer{
		Status: service.NewStatusService(db.Session, cfg.OnlineThreshold, db.Profile, db.LastSeen, realClock),
		Auth:   service.NewAuthService(db.Auth, db.Session, db.Profile),
		Ping:   service.NewPingService(db.Session, db.LastSeen),
	}
	handlers := di.HandlerContainer{
		Register:    transport.RegisterHandler(svcs.Auth),
		Login:       transport.LoginHandler(svcs.Auth),
		Ping:        transport.PingHandler(svcs.Ping),
		Status:      transport.StatusHandler(svcs.Status),
		BatchStatus: transport.BatchStatusHandler(svcs.Status),
	}
	server := transport.NewServer(handlers)
	err := server.Run(cfg.ServerPort)
	if err != nil {
		panic(err)
	}
}

func parseConfig() domain.Config {
	cfg := DefaultConfig
	flag.IntVar(&cfg.ServerPort, "p", 8080, "port to listen on")
	flag.DurationVar(&cfg.OnlineThreshold, "online-threshold", 10*time.Second, "threshold for determining if a user is online")
	flag.StringVar(&cfg.RedisAddress, "redis", "localhost:6379", "redis address in the format host:port")
	flag.Parse()

	envPort, set := os.LookupEnv("PORT")
	if set {
		if envPort, err := strconv.Atoi(envPort); err != nil {
			cfg.ServerPort = envPort
		}
	}
	envOnlineThresholdMillis, set := os.LookupEnv("ONLINE_THRESHOLD_MILLIS")
	if set {
		if envOnlineThresholdMillis, err := strconv.Atoi(envOnlineThresholdMillis); err != nil {
			cfg.OnlineThreshold = time.Duration(envOnlineThresholdMillis) * time.Millisecond
		}
	}
	envRedisAddress, set := os.LookupEnv("REDIS_ADDRESS")
	if set {
		cfg.RedisAddress = envRedisAddress
	}
	return cfg
}
