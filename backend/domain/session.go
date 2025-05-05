package domain

import (
	"errors"
	"time"
)

type Session struct {
	Token       string
	AccountId   string
	CreatedAt   time.Time
	RefreshedAt time.Time
}

var ErrSessionNotFound = errors.New("session not found")
var ErrInvalidSession = errors.New("invalid session")
var ErrSessionExpired = errors.New("session expired")

type SessionDao interface {
	Create(accountId string) Session
	Refresh(accountId string) Session
	GetBySessionToken(sessionToken string) (Session, bool)
	GetByAccountId(accountId string) (Session, bool)
	BatchGetByAccountId([]string) map[string]Session
}

type LastSeenDao interface {
	GetLastSeen(accountId string) (int64, error)
	SetLastSeen(accountId string, lastSeen int64) error
}
