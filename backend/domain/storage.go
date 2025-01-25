package domain

type Storage interface {
	Auth() AuthDao
}
