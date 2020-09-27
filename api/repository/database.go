package repository

import (
	"database/sql"
	"forum/config"
)

var DB *sql.DB

func InitDB() (err error) {
	if DB, err = sql.Open(config.DatabaseDriver, config.DatabasePath); err != nil {
		return err
	}
	DB.SetMaxIdleConns(100)
	if err = DB.Ping(); err != nil {
		return err
	}
	return nil
}

// CheckDBIntegrity creates necessary tables in the database, if they do not exist already
func CheckDBIntegrity() (err error) {
	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS users (
			id 	 	 	 INTEGER PRIMARY KEY,
			login     	 TEXT,
			password 	 BLOB,
			email 	 	 TEXT,
			display_name TEXT,
			created  	 TEXT,
			last_online  TEXT,
			session_id	 TEXT,
			role 		 INTEGER
		)`)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS posts (
			id				 INTEGER PRIMARY KEY,
			author_fkey		 INTEGER REFERENCES users(id),
			author_name_fkey TEXT REFERENCES users(display_name),
			title			 TEXT,
			content			 TEXT,
			created			 TEXT,
			updated			 TEXT,
			rating			 INTEGER
		)`)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS categories (
			id	 INTEGER PRIMARY KEY,
			name TEXT
		)`)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS posts_categories_bridge (
			id				 INTEGER PRIMARY KEY,
			post_id_fkey	 INTEGER REFERENCES posts(id),
			category_id_fkey INTEGER REFERENCES categories(id)
		)`)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS replies (
			id			 INTEGER PRIMARY KEY,
			author_fkey	 INTEGER REFERENCES users(id),
			content		 TEXT,
			created		 TEXT,
			post_id_fkey INTEGER REFERENCES posts(id)
		)`)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS reactions (
			id			 INTEGER PRIMARY KEY,
			post_id_fkey INTEGER REFERENCES posts(id),
			user_id_fkey INTEGER REFERENCES users(id),
			reaction	 INTEGER
		)`)
	if err != nil {
		return err
	}

	return nil
}
