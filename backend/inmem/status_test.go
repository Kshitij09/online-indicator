package inmem

import (
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestStatusCache_Update(t *testing.T) {
	cache := NewStatusCache()
	expected := domain.Status{Id: "1", IsOnline: true}
	cache.Update(expected)
	online, err := cache.IsOnline(expected.Id)
	if err != nil {
		t.Error(err)
	}
	if online != expected.IsOnline {
		t.Errorf("expected online=%v, got %v", expected.IsOnline, online)
	}
}

func TestStatusCache_LatestFetch(t *testing.T) {
	cache := NewStatusCache()
	expected := domain.Status{Id: "1", IsOnline: true}
	updated := domain.Status{Id: "2", IsOnline: false}
	cache.Update(expected)
	cache.Update(updated)
	updated.IsOnline = true
	cache.Update(updated)
	online, err := cache.IsOnline(updated.Id)
	if err != nil {
		t.Error(err)
	}
	if online != updated.IsOnline {
		t.Errorf("expected online=%v, got %v", expected.IsOnline, online)
	}
}

func TestStatusCache_ConcurrentReadWrite(t *testing.T) {
	cache := NewStatusCache()
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
			cache.Update(expectedEntries[i])
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
