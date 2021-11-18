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
	,[InsertAt]
	,[UpdateAt]
	,[InsertProfileID]
	,[InsertSystemID]
	,[UpdateProfileID]
	,[UpdateSystemID]
from
	[dbo].[Profiles]
`

const insertProfileQuery string = `
insert into [dbo].[Profiles] (
	 [Email]
	,[FirstName]
	,[MiddleName]
	,[LastName]
	,[Nickname]
	,[InsertProfileID]
	,[InsertSystemID]
	,[UpdateProfileID]
	,[UpdateSystemID]
	)
output [inserted].[ID]
values (
	 @p1
	,@p2
	,@p3
	,@p4
	,@p5
	,@p6
	,@p7
	,@p8
	,@p9
)
`

func GetProfiles(tx *sql.Tx, wps WherePhrases) ([]Profile, error) {
	profiles := make([]Profile, 0)
	query := getProfileQuery
	wp, values := wps.CreateWherePhrase(1)
	query += wp
	rows, err := tx.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		profile := Profile{}
		id := mssql.UniqueIdentifier{}
		insertProfileId := mssql.UniqueIdentifier{}
		insertSystemId := mssql.UniqueIdentifier{}
		updateProfileId := mssql.UniqueIdentifier{}
		updateSystemId := mssql.UniqueIdentifier{}
		if err := rows.Scan(
			&id,
			&profile.Email,
			&profile.FirstName,
			&profile.MiddleName,
			&profile.LastName,
			&profile.Nickname,
			&profile.ValidFlg,
			&profile.InsertAt,
			&profile.UpdateAt,
			&insertProfileId,
			&insertSystemId,
			&updateProfileId,
			&updateSystemId,
		); err != nil {
			return nil, err
		}
		profile.ID = id.String()
		profile.InsertProfileID = insertProfileId.String()
		profile.InsertSystemID = insertSystemId.String()
		profile.UpdateProfileID = updateProfileId.String()
		profile.UpdateSystemID = updateSystemId.String()
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
	insertProfileId := mssql.UniqueIdentifier{}
	insertSystemId := mssql.UniqueIdentifier{}
	updateProfileId := mssql.UniqueIdentifier{}
	updateSystemId := mssql.UniqueIdentifier{}
	insertProfileId.Scan(profile.InsertProfileID)
	insertSystemId.Scan(profile.InsertSystemID)
	updateProfileId.Scan(profile.UpdateProfileID)
	updateSystemId.Scan(profile.UpdateSystemID)
	err := tx.QueryRow(
		query,
		profile.Email,
		profile.FirstName,
		profile.MiddleName,
		profile.LastName,
		profile.Nickname,
		insertProfileId,
		insertSystemId,
		updateProfileId,
		updateSystemId).Scan(&id)
	if err != nil {
		return err
	}
	profile.ID = id.String()
	return nil
}

func (profile *Profile) update(tx *sql.Tx) error {
	up := NewUpdatePhrase("dbo", "Profile")
	up.ColumnValue["UpdateProfileId"] = profile.UpdateProfileID

	return nil
}
