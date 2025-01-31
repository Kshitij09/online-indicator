package inmem

import (
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	"sync"
)

type StatusCache struct {
	mu     sync.RWMutex
	online map[string]domain.Status
}

func NewStatusCache() *StatusCache {
	return &StatusCache{
		online: make(map[string]domain.Status),
	}
}

func (ctx *StatusCache) Update(status domain.Status) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.online[status.Id] = status
}

func (ctx *StatusCache) FetchAll(ids []string) []domain.Status {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	result := make([]domain.Status, 0)
	for _, id := range ids {
		isOnline := ctx.online[id].IsOnline
		result = append(result, domain.Status{Id: id, IsOnline: isOnline})
	}
	return result
}

func (ctx *StatusCache) IsOnline(id string) (bool, error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	status, ok := ctx.online[id]
	if ok {
		return status.IsOnline, nil
	} else {
		return false, fmt.Errorf("id %v not found", id)
	}
}
