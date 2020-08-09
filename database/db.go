package database

import (
	"database/sql"
	"forum/config"
	"log"
)

func Create(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE student (
		"idStudent" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"code" TEXT,
		"name" TEXT,
		"program" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create student table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("student table created")
}

func Connect() (*sql.DB, error) {
	db, err := sql.Open(config.DbDriver, config.DbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
