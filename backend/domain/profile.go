package domain

import "errors"

type Profile struct {
	UserId   string
	Username string
}

var ErrProfileAlreadyExists = errors.New("profile already exists")
var EmptyProfile = Profile{}

type ProfileDao interface {
	Create(Profile) error
	UsernameExists(string) bool
	GetByUserId(string) (Profile, bool)
	BatchGetByUserId([]string) map[string]Profile
}
