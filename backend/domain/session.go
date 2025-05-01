package domain

import (
	"errors"
	"time"
)

type Session struct {
	Id          string
	AccountId   string
	CreatedAt   time.Time
	RefreshedAt time.Time
}

var ErrSessionNotFound = errors.New("session not found")

type SessionDao interface {
	Create(accountId string) Session
	Refresh(sessionId string) bool
	GetBySessionId(sessionId string) (Session, bool)
	GetByAccountId(accountId string) (Session, bool)
	BatchGetByAccountId([]string) map[string]Session
}
