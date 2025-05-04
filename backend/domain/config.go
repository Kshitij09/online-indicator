package domain

import "time"

type Config struct {
	OnlineThreshold time.Duration
	ServerPort      int
	RedisAddress    string // host:port address
}
