package di

import (
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/service"
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
