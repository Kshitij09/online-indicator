package domain

import "errors"

type Account struct {
	Name  string
	Token string
}

var EmptyAccount = Account{}

var ErrAccountAlreadyExists = errors.New("account already exists")
var ErrAccountNotFound = errors.New("account not found")
var ErrEmptyName = errors.New("account name cannot be empty")

type AuthDao interface {
	Create(Account) (Account, error)
	Get(name string) (Account, bool)
	Update(Account) error
	Delete(name string) error
}
