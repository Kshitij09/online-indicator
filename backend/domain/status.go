package domain

import "time"

type Status struct {
	Id         string
	IsOnline   bool
	LastOnline time.Time
}

type StatusDao interface {
	UpdateOnline(id string, isOnline bool)
	IsOnline(id string) (bool, error)
	FetchAll(ids []string) []Status
}

type StatusService struct {
	status          StatusDao
	session         SessionDao
	onlineThreshold time.Duration
}

func NewStatusService(status StatusDao, session SessionDao, onlineThreshold time.Duration) StatusService {
	return StatusService{status: status, session: session, onlineThreshold: onlineThreshold}
}

func (ctx *StatusService) Ping(sessionId string) error {
	session, exists := ctx.session.GetBySessionId(sessionId)
	if !exists {
		return ErrSessionNotFound
	}
	ctx.status.UpdateOnline(session.Id, true)
	return nil
}
