package domain

import "github.com/google/uuid"

type TokenGenerator interface {
	Generate() string
}

type uuidTokenGenerator struct{}

func NewUUIDTokenGenerator() TokenGenerator {
	return uuidTokenGenerator{}
}

func (uuidTokenGenerator) Generate() string {
	return uuid.NewString()
}
