package inmem

import "github.com/Kshitij09/online-indicator/domain"

type Storage struct {
	Auth domain.RegisterDao
}

func NewStorage(tokenGen domain.TokenGenerator) *Storage {
	return &Storage{
		Auth: NewRegisterDao(tokenGen),
	}
}

func (ctx Storage) Register() domain.RegisterDao {
	return ctx.Auth
}
