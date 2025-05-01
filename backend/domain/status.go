package domain

import "time"

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
