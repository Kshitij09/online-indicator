package stubs

type StaticGenerator struct {
	StubValue string
}

func (ctx StaticGenerator) Generate() string {
	return ctx.StubValue
}
