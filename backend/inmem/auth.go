package inmem

import (
	"github.com/Kshitij09/online-indicator/domain"
)

type AuthCache struct {
	accounts        map[string]domain.Account
	accountNames    map[string]bool
	apiKeyGenerator domain.ApiKeyGenerator
	idGenerator     domain.IDGenerator
}

func NewAuthDao(apiKeyGenerator domain.ApiKeyGenerator, idGenerator domain.IDGenerator) domain.AuthDao {
	return &AuthCache{
		accounts:        make(map[string]domain.Account),
		accountNames:    make(map[string]bool),
		apiKeyGenerator: apiKeyGenerator,
		idGenerator:     idGenerator,
	}
}

func (ctx *AuthCache) Create(account domain.Account) (domain.Account, error) {
	if _, exists := ctx.accountNames[account.Name]; exists {
		return domain.EmptyAccount, domain.ErrAccountAlreadyExists
	}
	if account.Name == "" {
		return domain.EmptyAccount, domain.ErrEmptyName
	}
	account.ApiKey = ctx.apiKeyGenerator.Generate()
	account.Id = ctx.idGenerator.Generate()
	ctx.accounts[account.Id] = account
	ctx.accountNames[account.Name] = true
	return account, nil
}

func (ctx *AuthCache) Login(id string, apiKey string) (domain.Account, error) {
	acc, exists := ctx.accounts[id]
	if !exists {
		return domain.EmptyAccount, domain.ErrAccountNotFound
	}
	if acc.ApiKey != apiKey {
		return domain.EmptyAccount, domain.ErrInvalidCredentials
	}
	return acc, nil
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
