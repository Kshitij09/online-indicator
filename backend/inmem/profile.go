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

func NewProfileDao() domain.ProfileDao {
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
	if profile.Name == "" {
		return domain.ErrEmptyName
	}
	ctx.idLookup[profile.UserId] = &profile
	ctx.nameLookup[profile.Name] = &profile
	return nil
}

func (ctx *ProfileCache) NameExists(name string) bool {
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

func (ctx *ProfileCache) BatchGetByUserId(ids []string) map[string]domain.Profile {
	result := make(map[string]domain.Profile)
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	for _, id := range ids {
		profile, exists := ctx.idLookup[id]
		if exists {
			result[id] = *profile
		}
	}
	return result
}
