package service

import (
	"github.com/Kshitij09/online-indicator/domain"
	"time"
)

type StatusService struct {
	status          domain.StatusDao
	session         domain.SessionDao
	onlineThreshold time.Duration
}

func NewStatusService(status domain.StatusDao, session domain.SessionDao, onlineThreshold time.Duration) StatusService {
	return StatusService{status: status, session: session, onlineThreshold: onlineThreshold}
}

func (ctx *StatusService) Ping(sessionId string) error {
	session, exists := ctx.session.GetBySessionId(sessionId)
	if !exists {
		return domain.ErrSessionNotFound
	}
	ctx.status.UpdateOnline(session.Id, true)
	return nil
}
