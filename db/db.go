package db

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func MyNewSQLStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return db, nil
}
