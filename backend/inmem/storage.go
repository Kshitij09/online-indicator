package inmem

import (
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/jonboulle/clockwork"
)

type Storage struct {
	auth    domain.AuthDao
	session domain.SessionDao
	status  domain.StatusDao
	profile domain.ProfileDao
}

func NewStorage(
	tokenGen domain.TokenGenerator,
	sessionGen domain.SessionGenerator,
	clock clockwork.Clock,
	idGen domain.IDGenerator,
) *Storage {
	return &Storage{
		auth:    NewAuthDao(tokenGen, idGen),
		session: NewSessionCache(sessionGen, clock),
		status:  NewStatusCache(clock),
		profile: NewProfileCache(),
	}
}

func (ctx Storage) Auth() domain.AuthDao {
	return ctx.auth
}

func (ctx Storage) Session() domain.SessionDao { return ctx.session }

func (ctx Storage) Status() domain.StatusDao { return ctx.status }

func (ctx Storage) Profile() domain.ProfileDao { return ctx.profile }
