package models

import (
	"database/sql"
	"fmt"

	"github.com/atoyr/passed/database"
	"github.com/google/uuid"
)

type Authentication struct {
	Email    string `json:"email"`
	Key      string `json:"key"`
	Password string `json:"password"`
}

func (auth *Authentication) Signin(dbcontext *sql.DB) (bool, error) {
	tx, err := dbcontext.Begin()
	if err != nil {
		return false, err
	}

	profilewps := database.WherePhrases{}
	profilewps.Append(database.Equal, "Email", auth.Email)
	profiles, err := database.GetProfiles(tx, profilewps)
	if err != nil {
		return false, err
	}
	if len(profiles) != 1 {
		return false, fmt.Errorf("Profile not found : %d", len(profiles))
	}

	accountwps := database.WherePhrases{}
	accountwps.Append(database.Equal, "ProfileID", profiles[0].ID)
	accountwps.Append(database.Equal, "ValidFlg", true)
	accounts, err := database.GetAccounts(tx, accountwps)
	if err != nil {
		return false, err
	}
	if len(accounts) == 1 {
		u, err := uuid.Parse(auth.Key)
		if err != nil {
			return false, err
		}
		hash := GetSha3Hash([]byte(auth.Password))
		passhash := GetSha3Hash(hash, u.NodeID())

		email, err := AesDecript(passhash, accounts[0].Signature)
		if err != nil {
			return false, err
		}
		if auth.Email == string(email) {
			return true, nil
		}
	} else {
		return false, fmt.Errorf("account not found")
	}
	return false, fmt.Errorf("Login fail")
}
