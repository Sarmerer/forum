package database

import (
	"database/sql"
	"forum/config"
)

//Connect connects to the database
func Connect() (*sql.DB, error) {
	db, err := sql.Open(config.DbDriver, config.DbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
