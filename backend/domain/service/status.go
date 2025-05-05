package service

import (
	"errors"
	"github.com/Kshitij09/online-indicator/domain"
)

type StatusService struct {
	session  domain.SessionDao
	profile  domain.ProfileDao
	lastSeen domain.LastSeenDao
}

func NewStatusService(
	session domain.SessionDao,
	profile domain.ProfileDao,
	lastSeen domain.LastSeenDao,
) StatusService {
	return StatusService{
		session:  session,
		profile:  profile,
		lastSeen: lastSeen,
	}
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
		IsOnline:   ctx.isUserOnline(session.AccountId),
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
				IsOnline:   ctx.isUserOnline(session.AccountId),
				LastOnline: session.RefreshedAt,
			}
			merged[userId] = profileStatus
		}
	}
	return merged
}

func (ctx *StatusService) isUserOnline(accountId string) bool {
	_, err := ctx.lastSeen.GetLastSeen(accountId)
	return !errors.Is(err, domain.ErrSessionExpired)
}
