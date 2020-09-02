package database

import (
	"database/sql"
	"forum/config"
)

var DB *sql.DB

func InitDB() (err error) {
	if DB, err = sql.Open("sqlite3", config.DatabasePath); err != nil {
		return err
	}
	DB.SetMaxIdleConns(100)
	if err = DB.Ping(); err != nil {
		return err
	}
	return nil
}

// CheckIntegrity creates necessary tables in the database, if they do not exist already
func CheckIntegrity() (err error) {
	if err != nil {
		return
	}
	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS users (
			id 	 	 	INTEGER PRIMARY KEY,
			name     	TEXT,
			password 	BLOB,
			email 	 	TEXT,
			created  	TEXT,
			last_online TEXT,
			session_id	TEXT,
			role 		INTEGER
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS posts (
			id			INTEGER PRIMARY KEY,
			author_fkey	INTEGER,
			title		TEXT,
			content		TEXT,
			created		TEXT,
			updated		TEXT,
			rating		INTEGER,
			FOREIGN KEY(author_fkey) REFERENCES users(id)
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS categories (
			id	 INTEGER PRIMARY KEY,
			name TEXT
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS posts_categories_bridge (
			id				 INTEGER PRIMARY KEY,
			post_id_fkey	 INTEGER,
			category_id_fkey INTEGER,
			FOREIGN KEY(post_id_fkey) REFERENCES posts(id),
			FOREIGN KEY(category_id_fkey) REFERENCES categories(id)
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS "posts_reaction" (
			"post_reaction_id"	INTEGER,
			"post_id"	INTEGER,
			"user_id"	INTEGER,
			"reaction"	INTEGER,
			FOREIGN KEY("post_id") REFERENCES "posts"("post_id"),
			FOREIGN KEY("user_id") REFERENCES "users"("user_id"),
			PRIMARY KEY("post_reaction_id" AUTOINCREMENT)
		)`)
	if err != nil {
		return
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS replies (
			id			 INTEGER PRIMARY KEY,
			author		 INTEGER,
			content		 TEXT,
			created		 TEXT,
			post_id_fkey INTEGER,
			FOREIGN KEY(author) REFERENCES users(id),
			FOREIGN KEY(post_id_fkey) REFERENCES posts(id)
		)`)
	if err != nil {
		return
	}
	return
}
