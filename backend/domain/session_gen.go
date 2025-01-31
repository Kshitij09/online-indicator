package domain

import "github.com/google/uuid"

type SessionGenerator interface {
	Generate() string
}

type uuidSessionGenerator struct{}

func NewUUIDSessionGenerator() SessionGenerator {
	return uuidSessionGenerator{}
}

func (uuidSessionGenerator) Generate() string {
	return uuid.NewString()
}
