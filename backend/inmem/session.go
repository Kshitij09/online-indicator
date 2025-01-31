package inmem

import (
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/jonboulle/clockwork"
	"sync"
)

type SessionCache struct {
	mu              sync.RWMutex
	sessionIdLookup map[string]*domain.Session
	generator       domain.SessionGenerator
	clock           clockwork.Clock
}

func NewSessionCache(generator domain.SessionGenerator, clock clockwork.Clock) *SessionCache {
	return &SessionCache{
		sessionIdLookup: make(map[string]*domain.Session),
		generator:       generator,
		clock:           clock,
	}
}

func (ctx *SessionCache) Create(accountId string) domain.Session {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	session := domain.Session{
		Id:        ctx.generator.Generate(),
		AccountId: accountId,
		CreatedAt: ctx.clock.Now(),
	}
	ctx.sessionIdLookup[session.Id] = &session
	return session
}

func (ctx *SessionCache) GetBySessionId(sessionId string) (domain.Session, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	session, exists := ctx.sessionIdLookup[sessionId]
	if exists {
		return *session, true
	} else {
		return domain.Session{}, false
	}
}
