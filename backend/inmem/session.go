package inmem

import (
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/jonboulle/clockwork"
	"sync"
)

type SessionCache struct {
	mu              sync.RWMutex
	sessionIdLookup map[string]*domain.Session
	accountIdLookup map[string]*domain.Session
	generator       domain.SessionGenerator
	clock           clockwork.Clock
}

func NewSessionCache(generator domain.SessionGenerator, clock clockwork.Clock) *SessionCache {
	return &SessionCache{
		sessionIdLookup: make(map[string]*domain.Session),
		accountIdLookup: make(map[string]*domain.Session),
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
	ctx.accountIdLookup[accountId] = &session
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

func (ctx *SessionCache) GetByAccountId(accountId string) (domain.Session, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	session, exists := ctx.accountIdLookup[accountId]
	if exists {
		return *session, true
	} else {
		return domain.Session{}, false
	}
}

func (ctx *SessionCache) BatchGetByAccountId(ids []string) map[string]domain.Session {
	result := make(map[string]domain.Session)
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	for _, id := range ids {
		session, exists := ctx.accountIdLookup[id]
		if exists {
			result[id] = *session
		}
	}
	return result
}
