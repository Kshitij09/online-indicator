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
	status  StatusDao
	session SessionDao
}

func NewStatusService(status StatusDao, session SessionDao) StatusService {
	return StatusService{status: status, session: session}
}

func (ctx *StatusService) Ping(sessionId string) error {
	session, err := ctx.session.Get(sessionId)
	if err != nil {
		return err
	}
	ctx.status.UpdateOnline(session.Id, true)
	return nil
}
