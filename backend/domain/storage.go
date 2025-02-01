package domain

type Storage interface {
	Auth() AuthDao
	Session() SessionDao
	Status() StatusDao
	Profile() ProfileDao
}
