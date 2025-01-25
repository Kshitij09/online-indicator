package stubs

type StaticTokenGenerator struct {
	StubToken string
}

func (ctx StaticTokenGenerator) Generate() string {
	return ctx.StubToken
}
