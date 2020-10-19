package repository

import (
	"database/sql"
	"os"

	"github.com/sarmerer/forum/api/config"
)

var DB *sql.DB

func InitDB() (err error) {
	if _, err = os.Stat(config.DatabasePath); os.IsNotExist(err) {
		if err = os.Mkdir(config.DatabasePath, 0755); err != nil {
			return err
		}
	}
	if DB, err = sql.Open(config.DatabaseDriver, config.DatabasePath+"/"+config.DatabaseFileName); err != nil {
		return err
	}
	DB.SetMaxIdleConns(100)
	if err = DB.Ping(); err != nil {
		return err
	}
	return nil
}

// CheckDBIntegrity create database directory and tables, if they don't exist already
func CheckDBIntegrity() (err error) {
	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS users (
			id 	 	 	 INTEGER PRIMARY KEY,
			login     	 TEXT,
			password 	 BLOB,
			email 	 	 TEXT,
			avatar		 TEXT,
			display_name TEXT,
			created  	 TEXT,
			last_online  TEXT,
			session_id	 TEXT,
			role 		 INTEGER
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS posts (
			id				 INTEGER PRIMARY KEY,
			author_fkey		 INTEGER REFERENCES users(id),
			author_name_fkey TEXT REFERENCES users(display_name),
			title			 TEXT,
			content			 TEXT,
			created			 TEXT,
			updated			 TEXT
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS categories (
			id	 INTEGER PRIMARY KEY,
			name TEXT
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS posts_categories_bridge (
			id				 INTEGER PRIMARY KEY,
			post_id_fkey	 INTEGER REFERENCES posts(id),
			category_id_fkey INTEGER REFERENCES categories(id)
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS comments (
			id			 	 INTEGER PRIMARY KEY,
			author_fkey	 	 INTEGER REFERENCES users(id),
			author_name_fkey TEXT REFERENCES users(display_name),
			content		 	 TEXT,
			created		 	 TEXT,
			post_id_fkey 	 INTEGER REFERENCES posts(id),
			edited           INTEGER
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS reactions (
			id			 INTEGER PRIMARY KEY,
			post_id_fkey INTEGER REFERENCES posts(id),
			user_id_fkey INTEGER REFERENCES users(id),
			reaction	 INTEGER
		)`); err != nil {
		return err
	}

	return nil
}
