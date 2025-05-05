package service

import "github.com/Kshitij09/online-indicator/domain"

type PingService struct {
	session  domain.SessionDao
	lastSeen domain.LastSeenDao
}

func NewPingService(session domain.SessionDao, lastSeen domain.LastSeenDao) PingService {
	return PingService{
		session:  session,
		lastSeen: lastSeen,
	}
}

func (ctx *PingService) Ping(accountId, sessionToken string) error {
	session, exists := ctx.session.GetByAccountId(accountId)
	if !exists {
		return domain.ErrSessionNotFound
	}
	if session.Token != sessionToken {
		return domain.ErrInvalidSession
	}
	session = ctx.session.Refresh(session.AccountId)
	return ctx.lastSeen.SetLastSeen(accountId, session.RefreshedAt.Unix())
}
