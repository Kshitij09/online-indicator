package domain

import "errors"

type Account struct {
	Name  string
	Token string
	Id    string
}

var EmptyAccount = Account{}

var ErrAccountAlreadyExists = errors.New("account already exists")
var ErrAccountNotFound = errors.New("account not found")
var ErrEmptyName = errors.New("account name cannot be empty")
var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthDao interface {
	Create(Account) (Account, error)
	Login(name string, token string) error
	Update(Account) error
	Delete(name string) error
}

type AuthService struct {
	auth    AuthDao
	session SessionDao
}

func NewAuthService(auth AuthDao, session SessionDao) AuthService {
	return AuthService{
		auth:    auth,
		session: session,
	}
}

func (s AuthService) Login(name string, token string) (Session, error) {
	err := s.auth.Login(name, token)
	if err != nil {
		return Session{}, err
	}
	return s.session.Create(name), nil
}
