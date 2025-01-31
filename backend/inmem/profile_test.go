package inmem

import (
	"errors"
	"github.com/Kshitij09/online-indicator/domain"
	"reflect"
	"testing"
)

func TestProfileCache_Create(t *testing.T) {
	cache := NewProfileCache()
	expected := domain.Profile{Id: "1", Username: "test1"}
	err := cache.Create(expected)
	if err != nil {
		t.Error(err)
	}
	actual, exists := cache.GetByUserId(expected.Id)
	if !exists {
		t.Errorf("Profile not created")
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Created and Received profiles are different")
	}
}

func TestProfileCache_CreateDuplicate(t *testing.T) {
	cache := NewProfileCache()
	expected := domain.Profile{Id: "1", Username: "test1"}
	err := cache.Create(expected)
	if err != nil {
		t.Error(err)
	}
	err = cache.Create(expected)
	if !errors.Is(err, domain.ErrProfileAlreadyExists) {
		t.Errorf("expected error %v, got %v", domain.ErrProfileAlreadyExists, err)
	}
}

func TestProfileCache_UsernameExists(t *testing.T) {
	cache := NewProfileCache()
	expected := domain.Profile{Id: "1", Username: "test1"}
	if cache.UsernameExists(expected.Id) != false {
		t.Errorf("Profile should not exist initially")
	}
	err := cache.Create(expected)
	if err != nil {
		t.Error(err)
	}
	if cache.UsernameExists(expected.Username) != true {
		t.Errorf("Profile should exist after creation")
	}
}
