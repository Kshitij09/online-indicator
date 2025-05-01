package domain

import "time"

type Status struct {
	Id         string
	IsOnline   bool
	LastOnline time.Time
}

type ProfileStatus struct {
	Profile
	IsOnline   bool
	LastOnline time.Time
}

var EmptyProfileStatus = ProfileStatus{}

func OfflineProfileStatus(profile Profile, lastOnline time.Time) ProfileStatus {
	return ProfileStatus{
		Profile:    profile,
		IsOnline:   false,
		LastOnline: lastOnline,
	}
}

type StatusDao interface {
	UpdateOnline(id string, isOnline bool)
	Get(id string) (Status, error)
	BatchGet(ids []string) map[string]Status
}
