package database

import (
	"database/sql"
	"fmt"

	mssql "github.com/denisenkom/go-mssqldb"
)

type Account struct {
	ID        string `json:"id"`
	ProfileID string `json:"profile_id"`
	Signature []byte `json:"signature"`
	Private   []byte `json:"private"`
	Public    []byte `json:"public"`
	ValidFlg  bool   `json:"valid_flg"`
	History
}

const getAccountQuery string = `
select
	 [ID]
	,[ProfileID]
	,[Signature]
	,[Private]
	,[Public]
	,[ValidFlg]
	,[InsertDatetime]
	,[ModifiedDatetime]
	,[InsertAccountID]
	,[InsertSystemID]
	,[ModifiedAccountID]
	,[ModifiedSystemID]
from
	[dbo].[Accounts]
`

const insertAccountQuery string = `
insert into [dbo].[Accounts] (
	 [ProfileID]
	,[Signature]
	,[Private]
	,[Public]
	,[ValidFlg]
	,[InsertAccountID]
	,[InsertSystemID]
	,[ModifiedAccountID]
	,[ModifiedSystemID]
	)
output [inserted].[ID]
values (
	 @p1
	,@p2
	,@p3
	,@p4
	,@p5
	,[inserted].[ID]
	,@p6
	,[inserted].[ID]
	,@p7
)
`

func GetAuthenictions(tx *sql.Tx, wps WherePhrases) ([]Account, error) {
	accounts := make([]Account, 0)
	query := getAccountQuery
	wp, values := wps.CreateWherePhrase(1)
	query += wp
	rows, err := tx.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		account := Account{}
		id := mssql.UniqueIdentifier{}
		profileId := mssql.UniqueIdentifier{}
		insertAccountId := mssql.UniqueIdentifier{}
		insertSystemId := mssql.UniqueIdentifier{}
		modifiedAccountId := mssql.UniqueIdentifier{}
		modifiedSystemId := mssql.UniqueIdentifier{}
		if err := rows.Scan(
			id,
			profileId,
			account.Signature,
			account.Private,
			account.Public,
			account.ValidFlg,
			account.InsertDatetime,
			account.ModifiedDatetime,
			insertAccountId,
			insertSystemId,
			modifiedAccountId,
			modifiedSystemId,
		); err != nil {
			return nil, err
		}
		account.ID = id.String()
		account.ProfileID = profileId.String()
		account.InsertAccountID = insertAccountId.String()
		account.InsertSystemID = insertSystemId.String()
		account.ModifiedAccountID = modifiedAccountId.String()
		account.ModifiedSystemID = modifiedSystemId.String()
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
	id := mssql.UniqueIdentifier{}
	profileId := mssql.UniqueIdentifier{}
	insertSystemId := mssql.UniqueIdentifier{}
	modifiedSystemId := mssql.UniqueIdentifier{}
	profileId.Scan(account.ProfileID)
	insertSystemId.Scan(account.InsertSystemID)
	modifiedSystemId.Scan(account.ModifiedSystemID)
	err := tx.QueryRow(
		query,
		profileId,
		account.Signature,
		account.Private,
		account.Public,
		account.ValidFlg,
		insertSystemId,
		modifiedSystemId).Scan(&id)
	if err != nil {
		return err
	}
	account.ID = id.String()
	return nil
}

func (account *Account) update(tx *sql.Tx, beforeAccount Account) error {
	if beforeAccount.ModifiedDatetime.Equal(account.ModifiedDatetime) {

	} else {
		return fmt.Errorf("X0001")
	}

	return nil
}
