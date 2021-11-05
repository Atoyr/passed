package database

import "database/sql"

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
		if err := rows.Scan(
			profile.ID,
			profile.Email,
			profile.FirstName,
			profile.MiddleName,
			profile.LastName,
			profile.Nickname,
			profile.ValidFlg,
			profile.InsertDatetime,
			profile.ModifiedDatetime,
			profile.InsertAccountID,
			profile.InsertSystemID,
			profile.ModifiedAccountID,
			profile.ModifiedSystemID,
		); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}
