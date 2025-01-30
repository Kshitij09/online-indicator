package domain

type Status struct {
	Id       string
	IsOnline bool
}

type StatusDao interface {
	Update(Status)
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
	ctx.status.Update(Status{Id: session.Id, IsOnline: true})
	return nil
}
