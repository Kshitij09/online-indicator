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
	err := cache.Create(acc)
	if err != nil {
		t.Error(err)
	}
	created, exists := cache.Get(acc.Name)
	if !exists {
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
	err := cache.Create(acc)
	if !errors.Is(err, domain.ErrEmptyName) {
		t.Errorf("expected %s, got %s", domain.ErrEmptyName, err)
	}
}

func TestAuthCache_CreateExisting(t *testing.T) {
	tokenGen := StaticTokenGenerator{StubToken: "1"}
	cache := NewRegisterDao(tokenGen)
	acc := domain.Account{Name: "John Doe"}
	err := cache.Create(acc)
	if err != nil {
		t.Error(err)
	}
	err = cache.Create(acc)
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
	err := cache.Create(acc)
	if err != nil {
		t.Error(err)
	}
	_, exists = cache.Get(acc.Name)
	if !exists {
		t.Error("expected exists=true after creation, got false")
	}
}

type StaticTokenGenerator struct {
	StubToken string
}

func (ctx StaticTokenGenerator) Generate() string {
	return ctx.StubToken
}
