package service

import (
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/jonboulle/clockwork"
	"time"
)

type StatusService struct {
	session         domain.SessionDao
	profile         domain.ProfileDao
	onlineThreshold time.Duration
	clock           clockwork.Clock
}

func NewStatusService(
	session domain.SessionDao,
	onlineThreshold time.Duration,
	profile domain.ProfileDao,
	clock clockwork.Clock,
) StatusService {
	return StatusService{
		session:         session,
		onlineThreshold: onlineThreshold,
		profile:         profile,
		clock:           clock,
	}
}

func (ctx *StatusService) Ping(sessionId string) error {
	found := ctx.session.Refresh(sessionId)
	if !found {
		return domain.ErrSessionNotFound
	}
	return nil
}

func (ctx *StatusService) Status(accountId string) (domain.ProfileStatus, error) {
	profile, exists := ctx.profile.GetByUserId(accountId)
	if !exists {
		return domain.EmptyProfileStatus, domain.ErrAccountNotFound
	}
	session, exists := ctx.session.GetByAccountId(profile.UserId)
	if !exists {
		return domain.OfflineProfileStatus(profile, session.RefreshedAt), domain.ErrSessionNotFound
	}
	profileStatus := domain.ProfileStatus{
		Profile:    profile,
		IsOnline:   ctx.isSessionOnline(session.RefreshedAt),
		LastOnline: session.RefreshedAt,
	}
	return profileStatus, nil
}

func (ctx *StatusService) BatchStatus(ids []string) map[string]domain.ProfileStatus {
	profiles := ctx.profile.BatchGetByUserId(ids)
	sessions := ctx.session.BatchGetByAccountId(ids)
	merged := make(map[string]domain.ProfileStatus)
	for userId, profile := range profiles {
		session, exists := sessions[userId]
		if exists {
			profileStatus := domain.ProfileStatus{
				Profile:    profile,
				IsOnline:   ctx.isSessionOnline(session.RefreshedAt),
				LastOnline: session.RefreshedAt,
			}
			merged[userId] = profileStatus
		}
	}
	return merged
}

func (ctx *StatusService) isSessionOnline(lastRefresh time.Time) bool {
	return ctx.clock.Now().Sub(lastRefresh) <= ctx.onlineThreshold
}
