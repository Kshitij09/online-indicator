package service

import (
	"github.com/Kshitij09/online-indicator/domain"
	"time"
)

type StatusService struct {
	status          domain.StatusDao
	session         domain.SessionDao
	profile         domain.ProfileDao
	onlineThreshold time.Duration
}

func NewStatusService(
	status domain.StatusDao,
	session domain.SessionDao,
	onlineThreshold time.Duration,
	profile domain.ProfileDao,
) StatusService {
	return StatusService{
		status:          status,
		session:         session,
		onlineThreshold: onlineThreshold,
		profile:         profile,
	}
}

func (ctx *StatusService) Ping(sessionId string) error {
	session, exists := ctx.session.GetBySessionId(sessionId)
	if !exists {
		return domain.ErrSessionNotFound
	}
	ctx.status.UpdateOnline(session.AccountId, true)
	return nil
}

func (ctx *StatusService) Status(accountId string) (domain.ProfileStatus, error) {
	profile, exists := ctx.profile.GetByUserId(accountId)
	if !exists {
		return domain.EmptyProfileStatus, domain.ErrAccountNotFound
	}
	session, exists := ctx.session.GetByAccountId(profile.UserId)
	if !exists {
		return domain.OfflineProfileStatus(profile), domain.ErrSessionNotFound
	}
	status, err := ctx.status.Get(session.AccountId)
	if err != nil {
		return domain.EmptyProfileStatus, err
	}
	profileStatus := domain.ProfileStatus{
		Profile: profile,
		Status:  status,
	}
	return profileStatus, nil
}

func (ctx *StatusService) BatchStatus(ids []string) map[string]domain.ProfileStatus {
	profiles := ctx.profile.BatchGetByUserId(ids)
	statuses := ctx.status.BatchGet(ids)
	merged := make(map[string]domain.ProfileStatus)
	for userId, profile := range profiles {
		status, exists := statuses[userId]
		if exists {
			profileStatus := domain.ProfileStatus{
				Profile: profile,
				Status:  status,
			}
			merged[userId] = profileStatus
		}
	}
	return merged
}
