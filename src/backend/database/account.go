package database

import (
	"database/sql"
	"fmt"
)

type Account struct {
	ID         string `json:"id"`
	ProfileID  string `json:"profile_id"`
	Signature  []byte `json:"signature"`
	Private    []byte `json:"private"`
	Public     []byte `json:"public"`
	ValidFlg   bool   `json:"valid_flg"`
	UrgeSignin bool   `json:"urge_signin"`
	History
}

const getAccountQuery string = `
 select
 	ID
 	,ProfileID
 	,Signature
 	,Private
 	,Public
 	,ValidFlg
	,UrgeSignin
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
 	,Signature
 	,Private
 	,Public
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
			account.Signature,
			account.Private,
			account.Public,
			account.ValidFlg,
			account.UrgeSignin,
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

func (account *Account) Insert(tx *sql.Tx) error {
	wps := WherePhrases{}
	wps.Append(Equal, "ID", account.ID)
	accounts, err := GetAuthenictions(tx, wps)
	if err != nil {
		return err
	}
	if len(accounts) > 0 {
		return fmt.Errorf("Account Exists")
	} else {
		err = account.insert(tx)
	}
	if err != nil {
		return err
	}
	return nil
}

func (account *Account) Update(tx *sql.Tx) error {
	wps := WherePhrases{}
	wps.Append(Equal, "ID", account.ID)
	accounts, err := GetAuthenictions(tx, wps)
	if err != nil {
		return err
	}
	if len(accounts) > 0 {
		err = account.update(tx, accounts[0])
	} else {
		return fmt.Errorf("Account Not Found")
	}
	if err != nil {
		return err
	}
	return nil
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
		account.Signature,
		account.Private,
		account.Public,
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
