package test

import (
	"github.com/Kshitij09/online-indicator/domain"
	"time"
)

var Config = domain.Config{
	OnlineThreshold: 100 * time.Millisecond,
}
