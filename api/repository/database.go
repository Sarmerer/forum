package repository

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/utils"
)

// DB is a global database connection.
// Making so removes "connect to db" step from every operation with database,
// which reduces response time massively
var DB *sql.DB

// InitDB creates crucial for database and an entire API folders,
// and sets up a global database connection
func InitDB() (err error) {
	var (
		dbUser string = os.Getenv("DB_USER")
		dbPass string = os.Getenv("DB_PASS")
		dbAuth string
	)
	if err = utils.CreateFolderIfNotExists(config.DatabasePath); err != nil {
		return err
	}
	if err = utils.CreateFolderIfNotExists(config.DatabasePath + "/images"); err != nil {
		return err
	}

	if dbUser != "" && dbPass != "" {
		dbAuth = fmt.Sprintf("?_auth&_auth_user=%s&_auth_pass=%s", dbUser, dbPass)
	}

	if DB, err = sql.Open(config.DatabaseDriver, fmt.Sprintf("%s/%s%s", config.DatabasePath, config.DatabaseFileName, dbAuth)); err != nil {
		return err
	}
	DB.SetMaxIdleConns(100)
	if err = DB.Ping(); err != nil {
		return err
	}
	return nil
}

// CheckDBIntegrity creates tables, if they don't exist already
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
			is_image INTEGER,
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
