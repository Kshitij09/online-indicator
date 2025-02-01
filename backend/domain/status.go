package domain

import "time"

type Status struct {
	Id         string
	IsOnline   bool
	LastOnline time.Time
}

type ProfileStatus struct {
	Profile
	Status
}

var EmptyProfileStatus = ProfileStatus{}

func OfflineProfileStatus(profile Profile) ProfileStatus {
	return ProfileStatus{
		Profile: profile,
		Status: Status{
			Id:       profile.UserId,
			IsOnline: false,
		},
	}
}

type StatusDao interface {
	UpdateOnline(id string, isOnline bool)
	Get(id string) (Status, error)
	BatchGet(ids []string) map[string]Status
}
