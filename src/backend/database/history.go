package database

import "time"

type History struct {
	InsertAt        time.Time `json:"insert_at"`
	UpdateAt        time.Time `json:"update_at"`
	InsertAccountID string    `json:"insert_account_id"`
	InsertSystemID  string    `json:"insert_system_id"`
	UpdateAccountID string    `json:"update_account_id"`
	UpdateSystemID  string    `json:"update_system_id"`
}
