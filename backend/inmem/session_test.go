package inmem

import (
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/stubs"
	"github.com/jonboulle/clockwork"
	"maps"
	"reflect"
	"slices"
	"testing"
)

func TestSessionCache_CreateAndGet(t *testing.T) {
	staticGen := stubs.StaticGenerator{StubValue: "123"}
	clock := clockwork.NewFakeClock()
	cache := NewSessionCache(staticGen, clock)
	expectedSession := cache.Create("abc")
	session, exists := cache.GetBySessionToken("non-existent")
	if exists {
		t.Errorf("session '%s'should not exist", "non-existent")
	}
	session, exists = cache.GetBySessionToken(expectedSession.Token)
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
	if session.Token != expectedSession.Token {
		t.Errorf("session Token should be %s, got %s", expectedSession.Token, session.Token)
	}

	session, exists = cache.GetByAccountId(expectedSession.AccountId)
	if !exists {
		t.Errorf("session '%s'should exist", expectedSession.AccountId)
	}
	if session.CreatedAt != expectedCreationTime {
		t.Errorf("session createdAt should be %v, got %v", expectedCreationTime, session.CreatedAt)
	}
	if session.AccountId != "abc" {
		t.Errorf("session accountId should be %s, got %s", expectedSession.AccountId, session.AccountId)
	}
	if session.Token != expectedSession.Token {
		t.Errorf("session Token should be %s, got %s", expectedSession.Token, session.Token)
	}
}

func TestSessionCache_BatchGetByAccountId(t *testing.T) {
	staticGen := stubs.StaticGenerator{StubValue: "123"}
	clock := clockwork.NewFakeClock()
	cache := NewSessionCache(staticGen, clock)
	expected := make(map[string]domain.Session)
	for i := 0; i < 50; i++ {
		accountId := fmt.Sprintf("test%d", i)
		session := cache.Create(accountId)
		expected[session.AccountId] = session
	}
	actual := cache.BatchGetByAccountId(slices.Collect(maps.Keys(expected)))
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Created and Received profiles are different")
	}
}
