package di

import (
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/service"
	"github.com/Kshitij09/online-indicator/transport/handlers"
)

type DatabaseContainer struct {
	Auth     domain.AuthDao
	Session  domain.SessionDao
	Profile  domain.ProfileDao
	LastSeen domain.LastSeenDao
}

type ServiceContainer struct {
	Status service.StatusService
	Auth   service.AuthService
	Ping   service.PingService
}

type HandlerContainer struct {
	Register    handlers.Handler
	Login       handlers.Handler
	Status      handlers.Handler
	BatchStatus handlers.Handler
	Ping        handlers.Handler
}
