package domain

import (
	"errors"
	"time"
)

type Session struct {
	Id        string
	CreatedAt time.Time
}

var ErrInvalidSession = errors.New("invalid session")

type SessionDao interface {
	Create(Session)
	Get(id string) (Session, error)
}
