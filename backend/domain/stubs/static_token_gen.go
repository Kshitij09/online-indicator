package stubs

import (
	"github.com/Kshitij09/online-indicator/domain"
)

type StaticGenerator struct {
	StubValue string
}

func (ctx StaticGenerator) Generate() string {
	return ctx.StubValue
}

type StubLastSeenDao struct {
	lastSeen int64
	err      error
}

func (s *StubLastSeenDao) GetLastSeen(accountId string) (int64, error) {
	return s.lastSeen, s.err
}

func (s *StubLastSeenDao) SetLastSeen(accountId string, lastSeen int64) error {
	s.lastSeen = lastSeen
	s.err = nil
	return nil
}

func (s *StubLastSeenDao) SetAllOffline() {
	s.err = domain.ErrSessionExpired
}
