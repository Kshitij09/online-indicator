package domain

type Storage interface {
	Auth() AuthDao
	Session() SessionDao
	Profile() ProfileDao
}
