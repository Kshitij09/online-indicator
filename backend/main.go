package main

import (
	"flag"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/inmem"
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
	storage := inmem.NewStorage(apiKeyGen, sessionGen, realClock, idGen)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})
	server := transport.NewServer(storage, cfg, realClock, redisClient)
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
