package inmem

import (
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	"sync"
)

type StatusCache struct {
	mu     sync.RWMutex
	online map[string]bool
}

func NewStatusCache() *StatusCache {
	return &StatusCache{
		online: make(map[string]bool),
	}
}

func (ctx *StatusCache) Update(status domain.Status) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.online[status.Id] = status.IsOnline
}

func (ctx *StatusCache) FetchAll(ids []string) []domain.Status {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	result := make([]domain.Status, 0)
	for _, id := range ids {
		isOnline := ctx.online[id]
		result = append(result, domain.Status{Id: id, IsOnline: isOnline})
	}
	return result
}

func (ctx *StatusCache) IsOnline(id string) (bool, error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	online, ok := ctx.online[id]
	if ok {
		return online, nil
	} else {
		return false, fmt.Errorf("id %v not found", id)
	}
}
