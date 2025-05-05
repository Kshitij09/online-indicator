package inmem

import (
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/jonboulle/clockwork"
)

type Storage struct {
	auth    domain.AuthDao
	session domain.SessionDao
	profile domain.ProfileDao
}

func NewStorage(
	apiKeyGen domain.ApiKeyGenerator,
	sessionGen domain.SessionGenerator,
	clock clockwork.Clock,
	idGen domain.IDGenerator,
) *Storage {
	return &Storage{
		auth:    NewAuthDao(apiKeyGen, idGen),
		session: NewSessionDao(sessionGen, clock),
		profile: NewProfileDao(),
	}
}

func (ctx Storage) Auth() domain.AuthDao {
	return ctx.auth
}

func (ctx Storage) Session() domain.SessionDao { return ctx.session }

func (ctx Storage) Profile() domain.ProfileDao { return ctx.profile }
