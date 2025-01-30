package domain

type Status struct {
	Id       string
	IsOnline bool
}

type StatusDao interface {
	Update(Status)
	IsOnline(id string) (bool, error)
	FetchAll(ids []string) []Status
}
