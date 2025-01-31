package inmem

import (
	"errors"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/stubs"
	"testing"
)

func TestAuthCache_Create(t *testing.T) {
	tokenGen := stubs.StaticTokenGenerator{StubToken: "1"}
	cache := NewAuthDao(tokenGen)
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
	tokenGen := stubs.StaticTokenGenerator{StubToken: "1"}
	cache := NewAuthDao(tokenGen)
	acc := domain.Account{Name: ""}
	_, err := cache.Create(acc)
	if !errors.Is(err, domain.ErrEmptyName) {
		t.Errorf("expected %s, got %s", domain.ErrEmptyName, err)
	}
}

func TestAuthCache_CreateExisting(t *testing.T) {
	tokenGen := stubs.StaticTokenGenerator{StubToken: "1"}
	cache := NewAuthDao(tokenGen)
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
	tokenGen := stubs.StaticTokenGenerator{StubToken: "1"}
	cache := NewAuthDao(tokenGen)
	acc := domain.Account{Name: "John Doe"}
	err := cache.Login(acc.Name, tokenGen.StubToken)
	if err == nil {
		t.Errorf("expected %v initially, got nil", domain.ErrAccountNotFound)
	}
	_, err = cache.Create(acc)
	if err != nil {
		t.Error(err)
	}
	err = cache.Login(acc.Name, tokenGen.StubToken)
	if err != nil {
		t.Errorf("expected successful login, got %v", err)
	}
}
