package inmem

import (
	"github.com/Kshitij09/online-indicator/domain"
)

type AuthCache struct {
	accounts       map[string]domain.Account
	tokenGenerator domain.TokenGenerator
}

func NewAuthDao(tokenGenerator domain.TokenGenerator) domain.AuthDao {
	return &AuthCache{
		accounts:       make(map[string]domain.Account),
		tokenGenerator: tokenGenerator,
	}
}

func (ctx *AuthCache) Create(account domain.Account) (domain.Account, error) {
	if _, exists := ctx.accounts[account.Name]; exists {
		return domain.EmptyAccount, domain.ErrAccountAlreadyExists
	}
	if account.Name == "" {
		return domain.EmptyAccount, domain.ErrEmptyName
	}
	account.Token = ctx.tokenGenerator.Generate()
	ctx.accounts[account.Name] = account
	return account, nil
}

func (ctx *AuthCache) Get(name string) (domain.Account, bool) {
	acc, exists := ctx.accounts[name]
	if exists {
		return acc, true
	} else {
		return domain.Account{}, false
	}
}

func (ctx *AuthCache) Delete(name string) error {
	delete(ctx.accounts, name)
	return nil
}

func (ctx *AuthCache) Update(account domain.Account) error {
	if _, exists := ctx.accounts[account.Name]; !exists {
		return domain.ErrAccountNotFound
	}
	ctx.accounts[account.Name] = account
	return nil
}
