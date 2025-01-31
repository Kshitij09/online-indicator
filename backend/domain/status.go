package domain

import "time"

type Status struct {
	Id         string
	IsOnline   bool
	LastOnline time.Time
}

type StatusDao interface {
	UpdateOnline(id string, isOnline bool)
	Get(id string) (Status, error)
	FetchAll(ids []string) []Status
}
