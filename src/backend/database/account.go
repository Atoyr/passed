package database

import (
	"bytes"
	"database/sql"
	"fmt"

	mssql "github.com/denisenkom/go-mssqldb"
)

type Account struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Signature []byte `json:"signature"`
	Private   []byte `json:"private"`
	Public    []byte `json:"public"`
	ValidFlg  bool   `json:"valid_flg"`
	History
}

const getAccountQuery string = `
select
	 [ID]
	,[Email]
	,[Signature]
	,[Private]
	,[Public]
	,[ValidFlg]
	,[InsertAt]
	,[UpdateAt]
	,[InsertAccountID]
	,[InsertSystemID]
	,[UpdateAccountID]
	,[UpdateSystemID]
from
	[dbo].[Accounts]
`

const insertAccountQuery string = `
declare @nid uniqueidentifier = NEWID();
insert into [dbo].[Accounts] (
	 [ID]
	,[Email]
	,[Signature]
	,[Private]
	,[Public]
	,[InsertAccountID]
	,[InsertSystemID]
	,[UpdateAccountID]
	,[UpdateSystemID]
	)
output [inserted].[ID]
values (
	 @nid
	,@p1
	,@p2
	,@p3
	,@p4
	,@nid
	,@p5
	,@nid
	,@p6
);
`

func GetAccounts(tx *sql.Tx, wps WherePhrases) ([]Account, error) {
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
		insertProfileId := mssql.UniqueIdentifier{}
		insertSystemId := mssql.UniqueIdentifier{}
		updateProfileId := mssql.UniqueIdentifier{}
		updateSystemId := mssql.UniqueIdentifier{}
		if err := rows.Scan(
			&id,
			&account.Email,
			&account.Signature,
			&account.Private,
			&account.Public,
			&account.ValidFlg,
			&account.InsertAt,
			&account.UpdateAt,
			&insertProfileId,
			&insertSystemId,
			&updateProfileId,
			&updateSystemId,
		); err != nil {
			return nil, err
		}
		account.ID = id.String()
		account.InsertAccountID = insertProfileId.String()
		account.InsertSystemID = insertSystemId.String()
		account.UpdateAccountID = updateProfileId.String()
		account.UpdateSystemID = updateSystemId.String()
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (account *Account) Insert(tx *sql.Tx) error {
	wps := WherePhrases{}
	wps.Append(Equal, "ID", account.ID)
	accounts, err := GetAccounts(tx, wps)
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
	accounts, err := GetAccounts(tx, wps)
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

func (account *Account) insert(tx *sql.Tx) error {
	query := insertAccountQuery
	id := mssql.UniqueIdentifier{}
	insertSystemId := mssql.UniqueIdentifier{}
	updateSystemId := mssql.UniqueIdentifier{}
	insertSystemId.Scan(account.InsertSystemID)
	updateSystemId.Scan(account.UpdateSystemID)
	err := tx.QueryRow(
		query,
		account.Email,
		account.Signature,
		account.Private,
		account.Public,
		insertSystemId,
		updateSystemId).Scan(&id)
	if err != nil {
		return err
	}
	account.ID = id.String()
	return nil
}

func (account *Account) update(tx *sql.Tx, beforeAccount Account) error {
	if beforeAccount.UpdateAt.Equal(account.UpdateAt) {
		buffer := make([]byte, 4192)
		values := make([]interface{}, 0)
		isUpdate := false
		index := 1

		buffer = append(buffer, "UPDATE [dbo].[Accounts] "...)
		buffer = append(buffer, "SET "...)

		if !bytes.Equal(beforeAccount.Signature, account.Signature) {
			if isUpdate {
				buffer = append(buffer, " ,"...)
			}
			isUpdate = true

			buffer = append(buffer, []byte(fmt.Sprintf("Signature = @p%d", index))...)
			values = append(values, account.Signature)
			index++
		}
		if !bytes.Equal(beforeAccount.Private, account.Private) {
			isUpdate = true
			buffer = append(buffer, []byte(fmt.Sprintf("Private = @p%d, ", index))...)
			values = append(values, account.Private)
			index++
		}
		if !bytes.Equal(beforeAccount.Public, account.Public) {
			isUpdate = true
			buffer = append(buffer, []byte(fmt.Sprintf("Public = @p%d, ", index))...)
			values = append(values, account.Public)
			index++
		}
		if beforeAccount.ValidFlg != account.ValidFlg {
			isUpdate = true
			buffer = append(buffer, []byte(fmt.Sprintf("ValidFlg = @p%d, ", index))...)
			values = append(values, account.ValidFlg)
			index++
		}

		if !isUpdate {
			return fmt.Errorf("No Update Column")
		}

		buffer = append(buffer, "UpdateAt = GETDATE(), "...)
		buffer = append(buffer, []byte(fmt.Sprintf("UpdateAccountID = @p%d, ", index))...)
		values = append(values, account.UpdateAccountID)
		index++
		buffer = append(buffer, []byte(fmt.Sprintf("UpdateSystemID = @p%d ", index))...)
		values = append(values, account.UpdateSystemID)
		index++

		buffer = append(buffer, []byte(fmt.Sprintf("Where ID = @p%d ", index))...)
		values = append(values, account.ID)
		index++

		query := string(buffer)
		_, err := tx.Exec(query, values...)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("X0001")
	}

	return nil
}
