package inmem

import (
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/jonboulle/clockwork"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestStatusCache_Update(t *testing.T) {
	fakeClock := clockwork.NewFakeClock()
	expectedTime := fakeClock.Now()
	cache := NewStatusCache(fakeClock)
	expected := domain.Status{Id: "1", IsOnline: true, LastOnline: expectedTime}
	cache.UpdateOnline(expected.Id, expected.IsOnline)
	online, err := cache.IsOnline(expected.Id)
	if err != nil {
		t.Error(err)
	}
	if online != expected.IsOnline {
		t.Errorf("expected online=%v, got %v", expected.IsOnline, online)
	}
	lastOnline, err := cache.LastOnline(expected.Id)
	if err != nil {
		t.Error(err)
	}
	if lastOnline != expected.LastOnline {
		t.Errorf("expected lastOnline=%v, got %v", expected.LastOnline, lastOnline)
	}
}

func TestStatusCache_LatestFetch(t *testing.T) {
	fakeClock := clockwork.NewFakeClock()
	cache := NewStatusCache(fakeClock)
	updated := domain.Status{Id: "2", IsOnline: false}
	cache.UpdateOnline(updated.Id, updated.IsOnline)
	updated.IsOnline = true
	expectedOnlineTime := fakeClock.Now()
	cache.UpdateOnline(updated.Id, updated.IsOnline)
	online, err := cache.IsOnline(updated.Id)
	if err != nil {
		t.Error(err)
	}
	if online != updated.IsOnline {
		t.Errorf("expected online=%v, got %v", updated.IsOnline, online)
	}
	lastOnline, err := cache.LastOnline(updated.Id)
	if err != nil {
		t.Error(err)
	}
	if lastOnline != expectedOnlineTime {
		t.Errorf("expected lastOnline=%v, got %v", expectedOnlineTime, lastOnline)
	}
}

func TestStatusCache_ConcurrentReadWrite(t *testing.T) {
	cache := NewStatusCache(clockwork.NewFakeClock())
	expectedEntries := make([]domain.Status, 0, 1000)
	var wg sync.WaitGroup
	random := rand.New(rand.NewSource(40))
	for i := 0; i < 1000; i++ {
		isOnline := random.Intn(2) == 0
		expectedEntries = append(expectedEntries, domain.Status{Id: fmt.Sprint(i), IsOnline: isOnline})
		wg.Add(1)
		go func() {
			defer wg.Done()
			delay := time.Duration(random.Intn(11)) * time.Millisecond
			time.Sleep(delay)
			cache.UpdateOnline(expectedEntries[i].Id, expectedEntries[i].IsOnline)
		}()
	}
	wg.Wait()

	for i := 0; i < 1000; i++ {
		expected := expectedEntries[i]
		actualOnline, err := cache.IsOnline(expected.Id)
		if err != nil {
			t.Error(err)
		}
		if actualOnline != expected.IsOnline {
			t.Errorf("expected online for id %v=%v, got %v", expected.Id, expected.IsOnline, actualOnline)
		}
	}
}
