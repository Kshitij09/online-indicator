package service

import "github.com/Kshitij09/online-indicator/domain"

type AuthService struct {
	auth    domain.AuthDao
	session domain.SessionDao
	profile domain.ProfileDao
}

func NewAuthService(auth domain.AuthDao, session domain.SessionDao, profile domain.ProfileDao) AuthService {
	return AuthService{
		auth:    auth,
		session: session,
		profile: profile,
	}
}

func (s AuthService) Login(name string, token string) (domain.Session, error) {
	err := s.auth.Login(name, token)
	if err != nil {
		return domain.Session{}, err
	}
	return s.session.Create(name), nil
}

func (s AuthService) CreateAccount(account domain.Account) (domain.Account, error) {
	acc, err := s.auth.Create(account)
	if err != nil {
		return domain.EmptyAccount, err
	}
	profile := domain.Profile{
		UserId:   acc.Id,
		Username: acc.Name,
	}
	err = s.profile.Create(profile)
	if err != nil {
		return domain.EmptyAccount, err
	}
	return acc, nil
}
