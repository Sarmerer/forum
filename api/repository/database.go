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
			_id INTEGER PRIMARY KEY,
			username TEXT,
			password BLOB,
			email TEXT,
			avatar TEXT,
			alias TEXT,
			created INTEGER,
			last_active INTEGER,
			session_id TEXT,
			role INTEGER,
			verified INTEGER,
			oauth_provider TEXT
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS posts (
			_id INTEGER PRIMARY KEY,
			author_id_fkey INTEGER REFERENCES users(_id),
			title TEXT,
			content TEXT,
			created INTEGER,
			edited INTEGER,
			edit_reason TEXT
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS categories (_id INTEGER PRIMARY KEY, name TEXT)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS posts_categories_bridge (
			_id INTEGER PRIMARY KEY,
			post_id_fkey INTEGER REFERENCES posts(_id),
			category_id_fkey INTEGER REFERENCES categories(_id)
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS comments (
			_id INTEGER PRIMARY KEY,
			author_id_fkey INTEGER REFERENCES users(_id),
			post_id_fkey INTEGER REFERENCES posts(_id),
			parent_id_fkey INTEGER REFERENCES comments(_id),
			depth INTEGER,
			lineage TEXT,
			content TEXT,
			created INTEGER,
			edited INTEGER,
			deleted INTEGER
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS posts_reactions (
			_id			 INTEGER PRIMARY KEY,
			post_id_fkey INTEGER REFERENCES posts(_id),
			user_id_fkey INTEGER REFERENCES users(_id),
			reaction	 INTEGER
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS comments_reactions (
			_id			 INTEGER PRIMARY KEY,
			comment_id_fkey INTEGER REFERENCES comments(_id),
			user_id_fkey INTEGER REFERENCES users(_id),
			reaction	 INTEGER
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS user_saved_comments (
			_id			 INTEGER PRIMARY KEY,
			comment_id_fkey INTEGER REFERENCES comments(_id),
			user_id_fkey INTEGER REFERENCES users(_id)
		)`); err != nil {
		return err
	}

	if _, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS user_saved_posts (
			_id			 INTEGER PRIMARY KEY,
			post_id_fkey INTEGER REFERENCES comments(_id),
			user_id_fkey INTEGER REFERENCES users(_id)
		)`); err != nil {
		return err
	}
	return nil
}
