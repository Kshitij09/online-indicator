package inmem

import (
	"github.com/Kshitij09/online-indicator/domain"
	"sync"
)

type SessionCache struct {
	mu       sync.RWMutex
	sessions map[string]domain.Session
}

func NewSessionCache() *SessionCache {
	return &SessionCache{
		sessions: make(map[string]domain.Session),
	}
}

func (ctx *SessionCache) Create(session domain.Session) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.sessions[session.Id] = session
}

func (ctx *SessionCache) Get(id string) (domain.Session, error) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	session, exists := ctx.sessions[id]
	if exists {
		return session, nil
	} else {
		return domain.Session{}, domain.ErrInvalidSession
	}
}
