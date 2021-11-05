package database

import "time"

type History struct {
	InsertDatetime    time.Time `json:"insert_datetime"`
	ModifiedDatetime  time.Time `json:"modified_datetime"`
	InsertAccountID   string    `json:"insert_account_id"`
	InsertSystemID    string    `json:"insert_system_id"`
	ModifiedAccountID string    `json:"modified_account_id"`
	ModifiedSystemID  string    `json:"modified_system_id"`
}
