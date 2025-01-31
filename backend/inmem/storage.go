package inmem

import (
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/jonboulle/clockwork"
)

type Storage struct {
	auth    domain.AuthDao
	session domain.SessionDao
	status  domain.StatusDao
}

func NewStorage(tokenGen domain.TokenGenerator, clock clockwork.Clock) *Storage {
	return &Storage{
		auth:    NewAuthDao(tokenGen),
		session: NewSessionCache(),
		status:  NewStatusCache(clock),
	}
}

func (ctx Storage) Auth() domain.AuthDao {
	return ctx.auth
}

func (ctx Storage) Session() domain.SessionDao { return ctx.session }

func (ctx Storage) Status() domain.StatusDao { return ctx.status }
