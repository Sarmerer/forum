package database

import (
	"database/sql"
	"forum/config"
	"log"
)

//Connect connects to the database
func Connect() (*sql.DB, error) {
	db, err := sql.Open(config.DbDriver, config.DbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CheckIntegrity() (err error) {
	var db *sql.DB
	var execResult sql.Result
	var rowsAffected int64
	db, err = Connect()
	if err != nil {
		return err
	}
	execResult, err = db.Exec(`CREATE TABLE IF NOT EXISTS "users" (
		"user_id" INTEGER PRIMARY KEY,
		"user_name" TEXT,
		"user_password"	BLOB,
		"user_email" TEXT,
		"user_nickname"	TEXT,
		"user_created" TEXT,
		"user_last_online" TEXT
		"user_session_id" TEXT,
		"user_role" INTEGER)`,
	)
	if err != nil {
		return err
	}
	rowsAffected, err = execResult.RowsAffected()
	if rowsAffected > 0 {
		log.Println(`Recreated "users" table`)
	}
	if err != nil {
		return err
	}

	execResult, err = db.Exec(`CREATE TABLE IF NOT EXISTS "posts" (
		"post_id"	INTEGER,
		"post_by"	INTEGER,
		"post_category"	INTEGER,
		"post_name"	TEXT,
		"post_content"	TEXT,
		"post_created"	TEXT,
		"post_updated"	TEXT,
		"post_rating"	INTEGER,
		FOREIGN KEY("post_by") REFERENCES "users"("user_id"),
		FOREIGN KEY("post_category") REFERENCES "categories"("category_id"),
		PRIMARY KEY("post_id" AUTOINCREMENT))`,
	)
	if err != nil {
		return err
	}
	rowsAffected, err = execResult.RowsAffected()
	if rowsAffected > 0 {
		log.Println(`Recreated "posts" table`)
	}
	if err != nil {
		return err
	}

	execResult, err = db.Exec(`CREATE TABLE IF NOT EXISTS "categories" (
		"category_id"	INTEGER,
		"category_name"	TEXT UNIQUE,
		"category_description"	TEXT,
		PRIMARY KEY("category_id"))`)
	if err != nil {
		return err
	}
	rowsAffected, err = execResult.RowsAffected()
	if rowsAffected > 0 {
		log.Println(`Recreated "categories" table`)
	}
	if err != nil {
		return err
	}

	execResult, err = db.Exec(`CREATE TABLE IF NOT EXISTS "posts_reaction" (
		"post_reaction_id"	INTEGER,
		"post_id"	INTEGER,
		"user_id"	INTEGER,
		"reaction"	INTEGER,
		FOREIGN KEY("post_id") REFERENCES "posts"("post_id"),
		FOREIGN KEY("user_id") REFERENCES "users"("user_id"),
		PRIMARY KEY("post_reaction_id" AUTOINCREMENT))`,
	)
	if err != nil {
		return err
	}
	rowsAffected, err = execResult.RowsAffected()
	if rowsAffected > 0 {
		log.Println(`Recreated "reactions" table`)
	}
	if err != nil {
		return err
	}

	execResult, err = db.Exec(`CREATE TABLE IF NOT EXISTS "replies" (
		"reply_id"	INTEGER,
		"reply_content"	TEXT,
		"reply_date"	TEXT,
		"reply_post"	INTEGER,
		"reply_by"	INTEGER,
		FOREIGN KEY("reply_by") REFERENCES "users"("user_id"),
		FOREIGN KEY("reply_post") REFERENCES "posts"("post_id"),
		PRIMARY KEY("reply_id" AUTOINCREMENT))`,
	)
	if err != nil {
		return err
	}
	rowsAffected, err = execResult.RowsAffected()
	if rowsAffected > 0 {
		log.Println(`Recreated "replies" table`)
	}
	if err != nil {
		return err
	}
	return nil
}
