package domain

import (
	"errors"
	"time"
)

type Session struct {
	Id        string
	AccountId string
	CreatedAt time.Time
}

var ErrSessionNotFound = errors.New("session not found")

type SessionDao interface {
	Create(accountId string) Session
	GetBySessionId(sessionId string) (Session, bool)
}
