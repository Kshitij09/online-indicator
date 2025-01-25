package inmem

import (
	"github.com/Kshitij09/online-indicator/domain"
)

type AuthCache struct {
	accounts       map[string]domain.Account
	tokenGenerator domain.TokenGenerator
}

func NewAuthCache(tokenGenerator domain.TokenGenerator) *AuthCache {
	return &AuthCache{
		accounts:       make(map[string]domain.Account),
		tokenGenerator: tokenGenerator,
	}
}

func (ctx *AuthCache) Create(account domain.Account) error {
	if _, exists := ctx.accounts[account.Name]; exists {
		return domain.ErrAccountAlreadyExists
	}
	account.Token = ctx.tokenGenerator.Generate()
	ctx.accounts[account.Name] = account
	return nil
}

func (ctx *AuthCache) Get(name string) (domain.Account, bool) {
	acc, exists := ctx.accounts[name]
	if exists {
		return acc, true
	} else {
		return domain.Account{}, false
	}
}

func (ctx *AuthCache) Delete(name string) {
	delete(ctx.accounts, name)
}

func (ctx *AuthCache) Update(account domain.Account) error {
	if _, exists := ctx.accounts[account.Name]; !exists {
		return domain.ErrAccountNotFound
	}
	ctx.accounts[account.Name] = account
	return nil
}
