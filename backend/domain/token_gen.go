package domain

import "github.com/google/uuid"

type ApiKeyGenerator interface {
	Generate() string
}

type uuidApiKeyGenerator struct{}

func NewUUIDApiKeyGenerator() ApiKeyGenerator {
	return uuidApiKeyGenerator{}
}

func (uuidApiKeyGenerator) Generate() string {
	return uuid.NewString()
}
