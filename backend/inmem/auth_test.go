package inmem

import (
	"errors"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/stubs"
	"testing"
)

func TestAuthCache_Create(t *testing.T) {
	tokenGen := stubs.StaticGenerator{StubValue: "1"}
	cache := NewAuthDao(tokenGen, tokenGen)
	acc := domain.Account{Name: "John Doe"}
	created, err := cache.Create(acc)
	if err != nil {
		t.Error(err)
	}
	if created == domain.EmptyAccount {
		t.Errorf("account was not created")
	}
	if created.Token != tokenGen.StubValue {
		t.Errorf("token was not created")
	}
	if created.Id != tokenGen.StubValue {
		t.Errorf("id was not created")
	}
}

func TestAuthCache_Create_EmptyName(t *testing.T) {
	tokenGen := stubs.StaticGenerator{StubValue: "1"}
	cache := NewAuthDao(tokenGen, tokenGen)
	acc := domain.Account{Name: ""}
	_, err := cache.Create(acc)
	if !errors.Is(err, domain.ErrEmptyName) {
		t.Errorf("expected %s, got %s", domain.ErrEmptyName, err)
	}
}

func TestAuthCache_CreateExisting(t *testing.T) {
	tokenGen := stubs.StaticGenerator{StubValue: "1"}
	cache := NewAuthDao(tokenGen, tokenGen)
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
	tokenGen := stubs.StaticGenerator{StubValue: "1"}
	cache := NewAuthDao(tokenGen, tokenGen)
	acc := domain.Account{Name: "John Doe"}
	_, err := cache.Login(acc.Name, tokenGen.StubValue)
	if err == nil {
		t.Errorf("expected %v initially, got nil", domain.ErrAccountNotFound)
	}
	created, err := cache.Create(acc)
	if err != nil {
		t.Error(err)
	}
	acc, err = cache.Login(created.Id, tokenGen.StubValue)
	if err != nil {
		t.Errorf("expected successful login, got %v", err)
	}
}
