package database

import (
	"database/sql"
	"fmt"

	mssql "github.com/denisenkom/go-mssqldb"
)

type Profile struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Nickname   string `json:"nickname"`
	ValidFlg   bool   `json:"valid_flg"`
	History
}

const getProfileQuery string = `
select
	 [ID]
	,[Email]
	,[FirstName]
	,[MiddleName]
	,[LastName]
	,[Nickname]
	,[ValidFlg]
	,[InsertDatetime]
	,[ModifiedDatetime]
	,[InsertAccountID]
	,[InsertSystemID]
	,[ModifiedAccountID]
	,[ModifiedSystemID]
from
	[dbo].[Profile]
`

const insertProfileQuery string = `
insert into [dbo].[Profile] (
	 [Email]
	,[FirstName]
	,[MiddleName]
	,[LastName]
	,[Nickname]
	,[InsertAccountID]
	,[InsertSystemID]
	,[ModifiedAccountID]
	,[ModifiedSystemID]
	)
output [inserted].[ID]
value (
	,$1
	,$2
	,$3
	,$4
	,$5
	,$6
	,$7
	,$8
	,$9
)
`

func GetProfiles(tx *sql.Tx, wps WherePhrases) ([]Profile, error) {
	profiles := make([]Profile, 0)
	query := getAccountQuery
	wp, values := wps.CreateWherePhrase(1)
	query += wp
	rows, err := tx.Query(query, values)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		profile := Profile{}
		id := mssql.UniqueIdentifier{}
		insertAccountId := mssql.UniqueIdentifier{}
		insertSystemId := mssql.UniqueIdentifier{}
		modifiedAccountId := mssql.UniqueIdentifier{}
		modifiedSystemId := mssql.UniqueIdentifier{}
		if err := rows.Scan(
			id,
			profile.Email,
			profile.FirstName,
			profile.MiddleName,
			profile.LastName,
			profile.Nickname,
			profile.ValidFlg,
			profile.InsertDatetime,
			profile.ModifiedDatetime,
			insertAccountId,
			insertSystemId,
			modifiedAccountId,
			modifiedSystemId,
		); err != nil {
			return nil, err
		}
		profile.ID = id.String()
		profile.InsertAccountID = insertAccountId.String()
		profile.InsertSystemID = insertSystemId.String()
		profile.ModifiedAccountID = modifiedAccountId.String()
		profile.ModifiedSystemID = modifiedSystemId.String()
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (profile *Profile) Insert(tx *sql.Tx) error {
	wps := WherePhrases{}
	wps.Append(Equal, "ID", profile.ID)
	profiles, err := GetProfiles(tx, wps)
	if err != nil {
		return err
	}
	if len(profiles) > 0 {
		return fmt.Errorf("Account Exists")
	} else {
		err = profile.insert(tx)
	}
	if err != nil {
		return err
	}
	return nil
}

func (profile *Profile) insert(tx *sql.Tx) error {
	query := insertProfileQuery
	id := mssql.UniqueIdentifier{}
	insertAccountId := mssql.UniqueIdentifier{}
	insertSystemId := mssql.UniqueIdentifier{}
	modifiedAccountId := mssql.UniqueIdentifier{}
	modifiedSystemId := mssql.UniqueIdentifier{}
	insertAccountId.Scan(profile.InsertAccountID)
	insertSystemId.Scan(profile.InsertSystemID)
	modifiedAccountId.Scan(profile.ModifiedAccountID)
	modifiedSystemId.Scan(profile.ModifiedSystemID)
	err := tx.QueryRow(
		query,
		profile.Email,
		profile.FirstName,
		profile.MiddleName,
		profile.LastName,
		profile.Nickname,
		insertAccountId,
		insertSystemId,
		modifiedAccountId,
		modifiedSystemId).Scan(&id)
	if err != nil {
		return err
	}
	profile.ID = id.String()
	return nil
}
