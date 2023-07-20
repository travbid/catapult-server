package lib

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	return db, err
}
