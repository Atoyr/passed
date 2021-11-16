package orm

import (
	"database/sql"
)

type DB struct {
	db *sql.DB
	tx *sql.Tx
}

type Config struct {
	TableName string
	Scheme    string
}

func Open() (*DB, error) {

	return nil, nil
}

func (db *DB) Insert(data interface{}) {

}

func (db *DB) Update(where Where) {

}

func (db *DB) Delete() {

}

func (db *DB) Select() {

}

func (db *DB) Rollback() {

}

func (db *DB) Commit() {

}
