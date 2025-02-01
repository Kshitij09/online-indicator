package domain

import (
	"strconv"
	"sync/atomic"
)

type IDGenerator interface {
	Generate() string
}

type seqIdGenerator struct {
	counter atomic.Uint64
}

func NewSeqIdGenerator() IDGenerator {
	return &seqIdGenerator{}
}

func (id *seqIdGenerator) Generate() string {
	id.counter.Add(1)
	return strconv.FormatUint(id.counter.Load(), 10)
}
