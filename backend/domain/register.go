package domain

import "errors"

type Account struct {
	Name  string
	Token string
}

var ErrAccountAlreadyExists = errors.New("account already exists")
var ErrAccountNotFound = errors.New("account not found")

type RegisterDao interface {
	Create(Account) error
	Get(name string) (Account, bool)
	Update(Account) error
	Delete(name string) error
}
