package inmem

import (
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/jonboulle/clockwork"
	"sync"
	"time"
)

type StatusCache struct {
	mu     sync.RWMutex
	online map[string]*domain.Status
	clock  clockwork.Clock
}

func NewStatusCache(clock clockwork.Clock) *StatusCache {
	return &StatusCache{
		online: make(map[string]*domain.Status),
		clock:  clock,
	}
}

func (ctx *StatusCache) UpdateOnline(id string, isOnline bool) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	status, exists := ctx.online[id]
	if !exists {
		status = &domain.Status{
			Id:       id,
			IsOnline: isOnline,
		}
	}
	if isOnline {
		status.LastOnline = ctx.clock.Now()
	}
	status.IsOnline = isOnline
	ctx.online[id] = status
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

func (ctx *StatusCache) LastOnline(id string) (time.Time, error) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	if status, ok := ctx.online[id]; ok {
		return status.LastOnline, nil
	} else {
		return time.Time{}, fmt.Errorf("id %v not found", id)
	}
}
