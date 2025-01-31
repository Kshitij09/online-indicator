package inmem

import (
	"github.com/Kshitij09/online-indicator/domain/stubs"
	"github.com/jonboulle/clockwork"
	"testing"
)

func TestSessionCache_CreateAndGet(t *testing.T) {
	staticGen := stubs.StaticGenerator{StubValue: "123"}
	clock := clockwork.NewFakeClock()
	cache := NewSessionCache(staticGen, clock)
	expectedSession := cache.Create("abc")
	session, exists := cache.Get("non-existent")
	if exists {
		t.Errorf("session '%s'should not exist", "non-existent")
	}
	session, exists = cache.Get(expectedSession.Id)
	if !exists {
		t.Errorf("session '%s'should exist", expectedSession.AccountId)
	}
	expectedCreationTime := clock.Now()
	if session.CreatedAt != expectedCreationTime {
		t.Errorf("session createdAt should be %v, got %v", expectedCreationTime, session.CreatedAt)
	}
	if session.AccountId != "abc" {
		t.Errorf("session accountId should be %s, got %s", expectedSession.AccountId, session.AccountId)
	}
	if session.Id != expectedSession.Id {
		t.Errorf("session Id should be %s, got %s", expectedSession.Id, session.Id)
	}
}
