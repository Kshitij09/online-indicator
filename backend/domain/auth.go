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

type LoginService struct {
	auth    AuthDao
	session SessionDao
}

func NewLoginService(auth AuthDao, session SessionDao) LoginService {
	return LoginService{
		auth:    auth,
		session: session,
	}
}

func (s LoginService) Login(name string, token string) (Session, error) {
	err := s.auth.Login(name, token)
	if err != nil {
		return Session{}, err
	}
	return s.session.Create(name), nil
}
