package inmem

import (
	"errors"
	"github.com/Kshitij09/online-indicator/domain"
	"testing"
)

func TestAuthCache_Create(t *testing.T) {
	tokenGen := StaticTokenGenerator{StubToken: "1"}
	cache := NewRegisterDao(tokenGen)
	acc := domain.Account{Name: "John Doe"}
	created, err := cache.Create(acc)
	if err != nil {
		t.Error(err)
	}
	if created == domain.EmptyAccount {
		t.Errorf("account was not created")
	}
	if created.Token != tokenGen.StubToken {
		t.Errorf("token was not created")
	}
}

func TestAuthCache_Create_EmptyName(t *testing.T) {
	tokenGen := StaticTokenGenerator{StubToken: "1"}
	cache := NewRegisterDao(tokenGen)
	acc := domain.Account{Name: ""}
	_, err := cache.Create(acc)
	if !errors.Is(err, domain.ErrEmptyName) {
		t.Errorf("expected %s, got %s", domain.ErrEmptyName, err)
	}
}

func TestAuthCache_CreateExisting(t *testing.T) {
	tokenGen := StaticTokenGenerator{StubToken: "1"}
	cache := NewRegisterDao(tokenGen)
	acc := domain.Account{Name: "John Doe"}
	_, err := cache.Create(acc)
	if err != nil {
		t.Error(err)
	}
	_, err = cache.Create(acc)
	if !errors.Is(err, domain.ErrAccountAlreadyExists) {
		t.Errorf("expected %s, got %s", domain.ErrAccountAlreadyExists, err)
	}
}

func TestAuthCache_Get(t *testing.T) {
	tokenGen := StaticTokenGenerator{StubToken: "1"}
	cache := NewRegisterDao(tokenGen)
	acc := domain.Account{Name: "John Doe"}
	_, exists := cache.Get(acc.Name)
	if exists {
		t.Error("expected exists=false initially, got true")
	}
	created, err := cache.Create(acc)
	if err != nil {
		t.Error(err)
	}
	fetched, exists := cache.Get(acc.Name)
	if !exists {
		t.Error("expected exists=true after creation, got false")
	}
	if fetched != created {
		t.Errorf("expected %s, got %s", acc, fetched)
	}
}

type StaticTokenGenerator struct {
	StubToken string
}

func (ctx StaticTokenGenerator) Generate() string {
	return ctx.StubToken
}
