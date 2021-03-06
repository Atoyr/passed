package models

import (
	"database/sql"
	"time"

	"github.com/atoyr/passed/database"
)

type User struct {
	AccountID                string    `jso:"account_id"`
	ProfileID                string    `jso:"profile_id"`
	Signature                string    `json:"signatu"`
	Private                  string    `json:"private"`
	Public                   string    `json:"public"`
	Email                    string    `json:"email"`
	FirstName                string    `json:"first_name"`
	MiddleName               string    `json:"middle_name"`
	LastName                 string    `json:"last_name"`
	Nickname                 string    `json:"nickname"`
	AccountValidFlg          bool      `json:"account_valid_flg"`
	ProfileValidFlg          bool      `json:"profile_valid_flg"`
	UrgeSignin               bool      `json:"urge_signin"`
	AccountInsertDatetime    time.Time `json:"account_insert_datetime"`
	AccountModifiedDatetime  time.Time `json:"account_modified_datetime"`
	AccountInsertAccountID   string    `json:"account_insert_account_id"`
	AccountInsertSystemID    string    `json:"account_insert_system_id"`
	AccountModifiedAccountID string    `json:"account_modified_account_id"`
	AccountModifiedSystemID  string    `json:"account_modified_system_id"`
	ProfileInsertDatetime    time.Time `json:"profile_insert_datetime"`
	ProfileModifiedDatetime  time.Time `json:"profile_modified_datetime"`
	ProfileInsertAccountID   string    `json:"profile_insert_account_id"`
	ProfileInsertSystemID    string    `json:"profile_insert_system_id"`
	ProfileModifiedAccountID string    `json:"profile_modified_account_id"`
	ProfileModifiedSystemID  string    `json:"profile_modified_system_id"`
}

func (user *User) Create(dbcontext *sql.DB) error {
	tx, err := dbcontext.Begin()
	if err != nil {
		return err
	}
	wps := database.WherePhrases{}
	wps.Append(database.Equal, "Email", user.Email)
	database.GetProfiles(tx, wps)
	user.AccountID = ""
	user.ProfileID = ""

	return nil
}
