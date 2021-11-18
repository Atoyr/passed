package database

import "time"

type History struct {
	InsertAt        time.Time `json:"insert_at"`
	UpdateAt        time.Time `json:"update_at"`
	InsertProfileID string    `json:"insert_profile_id"`
	InsertSystemID  string    `json:"insert_system_id"`
	UpdateProfileID string    `json:"update_profile_id"`
	UpdateSystemID  string    `json:"update_system_id"`
}
