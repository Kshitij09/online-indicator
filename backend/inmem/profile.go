package inmem

import (
	"github.com/Kshitij09/online-indicator/domain"
	"sync"
)

type ProfileCache struct {
	mu         sync.RWMutex
	idLookup   map[string]*domain.Profile
	nameLookup map[string]*domain.Profile
}

func NewProfileCache() domain.ProfileDao {
	return &ProfileCache{
		idLookup:   make(map[string]*domain.Profile),
		nameLookup: make(map[string]*domain.Profile),
	}
}

func (ctx *ProfileCache) Create(profile domain.Profile) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if _, exists := ctx.idLookup[profile.UserId]; exists {
		return domain.ErrProfileAlreadyExists
	}
	if profile.Username == "" {
		return domain.ErrEmptyName
	}
	ctx.idLookup[profile.UserId] = &profile
	ctx.nameLookup[profile.Username] = &profile
	return nil
}

func (ctx *ProfileCache) UsernameExists(name string) bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	_, exists := ctx.nameLookup[name]
	return exists
}

func (ctx *ProfileCache) GetByUserId(id string) (domain.Profile, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	profile, exists := ctx.idLookup[id]
	if exists {
		return *profile, exists
	} else {
		return domain.EmptyProfile, false
	}
}
