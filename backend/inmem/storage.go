package inmem

import "github.com/Kshitij09/online-indicator/domain"

type Storage struct {
	auth domain.AuthDao
}

func NewStorage(tokenGen domain.TokenGenerator) *Storage {
	return &Storage{
		auth: NewAuthDao(tokenGen),
	}
}

func (ctx Storage) Auth() domain.AuthDao {
	return ctx.auth
}
