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

func (s AuthService) Login(id string, token string) (domain.Session, error) {
	acc, err := s.auth.Login(id, token)
	if err != nil {
		return domain.Session{}, err
	}
	return s.session.Create(acc.Id), nil
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
