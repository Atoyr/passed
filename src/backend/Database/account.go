package database

import (
	"database/sql"
	"fmt"
)

type Account struct {
	ID        string `json:"id"`
	ProfileID string `json:"profile_id"`
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Shared    string `json:"shared"`
	ValidFlg  bool   `json:"valid_flg"`
	History
}

const getAccountQuery string = `
 select
 	ID
 	,ProfileID
 	,Primary
 	,Secondary
 	,Shared
 	,ValidFlg
 	,InsertDatetime
 	,ModifiedDatetime
 	,InsertAccountID
 	,InsertSystemID
 	,ModifiedAccountID
 	,ModifiedSystemID
 from
 	dbo.Accounts
 `

const insertAccountQuery string = `
insert into dbo.Accounts (
 	 ProfileID
 	,Primary
 	,Secondary
 	,Shared
 	,ValidFlg
 	,InsertAccountID
 	,InsertSystemID
 	,ModifiedAccountID
 	,ModifiedSystemID
	)
	output inserted.ID
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

func GetAuthenictions(tx *sql.Tx, wps WherePhrases) ([]Account, error) {
	accounts := make([]Account, 0)
	query := getAccountQuery
	wp, values := wps.CreateWherePhrase(1)
	query += wp
	rows, err := tx.Query(query, values)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		account := Account{}
		if err := rows.Scan(
			account.ID,
			account.ProfileID,
			account.Primary,
			account.Secondary,
			account.Shared,
			account.ValidFlg,
			account.InsertDatetime,
			account.ModifiedDatetime,
			account.InsertAccountID,
			account.InsertSystemID,
			account.ModifiedAccountID,
			account.ModifiedSystemID,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (account *Account) InsertOrUpdate(tx *sql.Tx) error {
	wps := WherePhrases{}
	wps.Append(Equal, "ID", account.ID)
	accounts, err := GetAuthenictions(tx, wps)
	if err != nil {
		return err
	}
	if len(accounts) > 0 {
		err = account.update(tx, accounts[0])
	} else {
		err = account.insert(tx)
	}
	if err != nil {
		return err
	}
	return nil
}

func (account *Account) insert(tx *sql.Tx) error {
	query := insertAccountQuery
	id := ""
	err := tx.QueryRow(
		query,
		account.ProfileID,
		account.Primary,
		account.Secondary,
		account.Shared,
		account.ValidFlg,
		account.InsertAccountID,
		account.InsertSystemID,
		account.ModifiedAccountID,
		account.ModifiedSystemID).Scan(&id)
	if err != nil {
		return err
	}
	account.ID = id
	return nil
}

func (account *Account) update(tx *sql.Tx, beforeAccount Account) error {
	if beforeAccount.ModifiedDatetime.Equal(account.ModifiedDatetime) {

	} else {
		return fmt.Errorf("X0001")
	}

	return nil
}
