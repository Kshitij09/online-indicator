package inmem

import (
	"container/list"
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/jonboulle/clockwork"
	"sync"
	"time"
)

type StatusCache struct {
	mu    sync.RWMutex
	items map[string]*list.Element
	list  *list.List
	clock clockwork.Clock
}

func NewStatusCache(clock clockwork.Clock) *StatusCache {
	return &StatusCache{
		items: make(map[string]*list.Element),
		clock: clock,
		list:  list.New(),
	}
}

func (ctx *StatusCache) UpdateOnline(id string, isOnline bool) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	element, exists := ctx.items[id]
	if !exists {
		newItem := &domain.Status{
			Id:       id,
			IsOnline: isOnline,
		}
		element = ctx.list.PushBack(newItem)
		ctx.items[id] = element
	}
	status := element.Value.(*domain.Status)
	if isOnline {
		status.LastOnline = ctx.clock.Now()
		ctx.list.MoveToFront(element)
	}
	status.IsOnline = isOnline
}

func (ctx *StatusCache) FetchAll(ids []string) []domain.Status {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	result := make([]domain.Status, 0)
	for _, id := range ids {
		status, err := ctx.get(id)
		if err != nil {
			continue
		}
		result = append(result, *status)
	}
	return result
}

func (ctx *StatusCache) Get(id string) (domain.Status, error) {
	status, err := ctx.get(id)
	if err != nil {
		return domain.Status{}, err
	} else {
		return *status, nil
	}
}

func (ctx *StatusCache) LastOnline(id string) (time.Time, error) {
	status, err := ctx.get(id)
	if err != nil {
		return time.Time{}, err
	} else {
		return status.LastOnline, nil
	}
}

func (ctx *StatusCache) get(id string) (*domain.Status, error) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	if element, ok := ctx.items[id]; ok {
		return element.Value.(*domain.Status), nil
	} else {
		return nil, fmt.Errorf("id %v not found", id)
	}
}
