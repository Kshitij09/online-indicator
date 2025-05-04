package stubs

type StaticGenerator struct {
	StubValue string
}

func (ctx StaticGenerator) Generate() string {
	return ctx.StubValue
}

type StubLastSeenDao struct {
	lastSeen int64
}

func (s StubLastSeenDao) GetLastSeen(accountId string) (int64, error) {
	return s.lastSeen, nil
}

func (s StubLastSeenDao) SetLastSeen(accountId string, lastSeen int64) error {
	s.lastSeen = lastSeen
	return nil
}
