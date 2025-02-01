package inmem

import (
	"errors"
	"github.com/Kshitij09/online-indicator/domain"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"testing"
)

func TestProfileCache_Create(t *testing.T) {
	cache := NewProfileCache()
	expected := domain.Profile{UserId: "1", Username: "test1"}
	err := cache.Create(expected)
	if err != nil {
		t.Error(err)
	}
	actual, exists := cache.GetByUserId(expected.UserId)
	if !exists {
		t.Errorf("Profile not created")
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Created and Received profiles are different")
	}
}

func TestProfileCache_CreateDuplicate(t *testing.T) {
	cache := NewProfileCache()
	expected := domain.Profile{UserId: "1", Username: "test1"}
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
	expected := domain.Profile{UserId: "1", Username: "test1"}
	if cache.UsernameExists(expected.UserId) != false {
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

func TestProfileCache_BatchGetByUserId(t *testing.T) {
	cache := NewProfileCache()
	expected := make(map[string]domain.Profile)
	for i := 0; i < 50; i++ {
		userId := strconv.Itoa(i)
		userName := "test" + userId
		profile := domain.Profile{UserId: userId, Username: userName}
		expected[userId] = profile
		err := cache.Create(profile)
		if err != nil {
			t.Error(err)
		}
	}
	actual := cache.BatchGetByUserId(slices.Collect(maps.Keys(expected)))
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Created and Received profiles are different")
	}
}
