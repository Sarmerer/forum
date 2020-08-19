package models

import (
	"database/sql"
	"errors"
	"forum/api/entities"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const timeLayout = "2006-01-02 15:04:05"

//UserModel helps performing CRUD operations
type UserModel struct {
	DB *sql.DB
}

//NewUserModel creates an instance of UserModel
func NewUserModel(db *sql.DB) (*UserModel, error) {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (user_id INTEGER PRIMARY KEY, user_name TEXT, user_password	BLOB, user_email TEXT, user_nickname	TEXT, user_created	TEXT, user_last_online	TEXT, user_session_id TEXT, user_role INTEGER)")
	if err != nil {
		return nil, err
	}
	statement.Exec()
	return &UserModel{db}, nil
}

//FindAll returns all users in the database
func (um *UserModel) FindAll() ([]entities.User, error) {
	rows, queryErr := um.DB.Query("SELECT * FROM users")
	if queryErr != nil {
		return nil, queryErr
	}
	var users []entities.User

	for rows.Next() {
		var user entities.User
		var created, lastOnline string
		rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
		date, dateErr := time.Parse(timeLayout, created)
		if dateErr != nil {
			return users, dateErr
		}
		user.Created = date
		date, dateErr = time.Parse(timeLayout, lastOnline)
		if dateErr != nil {
			return users, dateErr
		}
		user.LastOnline = date
		users = append(users, user)
	}
	return users, nil
}

//Find returns a specific user from the database
func (um *UserModel) Find(id int64) (entities.User, error) {
	var user entities.User
	rows, err := um.DB.Query("SELECT * FROM users WHERE user_id = ?", id)
	if err != nil {
		return user, err
	}
	for rows.Next() {
		var created, lastOnline string
		rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
		date, err := time.Parse(timeLayout, created)
		if err != nil {
			return user, err
		}
		user.Created = date
		date, err = time.Parse(timeLayout, lastOnline)
		if err != nil {
			return user, err
		}
		user.LastOnline = date
	}
	return user, nil
}

//Create adds a new user to the database
func (um *UserModel) Create(user *entities.User) error {
	statement, err := um.DB.Prepare("INSERT INTO users (user_name, user_password, user_email, user_nickname, user_created, user_last_online, user_session_id, user_role) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("internal server error")
	}
	res, err := statement.Exec(user.Name, user.Password, user.Email, user.Nickname, time.Now().Format(timeLayout), time.Now().Format(timeLayout), user.SessionID, user.Role)
	if err.Error() == "UNIQUE constraint failed: users.user_email" {
		return errors.New("email is not unique")
	} else if err.Error() == "UNIQUE constraint failed: users.user_name" {
		return errors.New("login is not unique")
	} else if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.New("internal server error")
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("internal server error")
}

//Delete deletes user from the database
func (um *UserModel) Delete(id int64) bool {
	res, err := um.DB.Exec("DELETE FROM users WHERE user_id = ?", id)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}

//Update updates existing user in the database
func (um *UserModel) Update(user *entities.User) bool {
	statement, err := um.DB.Prepare("UPDATE users SET user_name = ?, user_password = ?, user_email = ?, user_nickname = ?, user_created = ?, user_last_online = ?, user_session_id = ?, user_role = ? WHERE user_id = ?")
	if err != nil {
		return false
	}
	res, err := statement.Exec(user.Name, user.Password, user.Email, user.Nickname, user.Created.Format(timeLayout), user.LastOnline.Format(timeLayout), user.SessionID, user.Role, user.ID)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}

//Validate checks if the user is logged in using session id
func (um *UserModel) Validate(id string) (bool, error) {
	err := um.DB.QueryRow("SELECT user_name FROM users WHERE user_session_id = ?", id).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

//FindByNameOrEmail finds a user by name or email
func (um *UserModel) FindByNameOrEmail(login string) (entities.User, error) {
	var user entities.User
	rows, err := um.DB.Query("SELECT * FROM users WHERE user_name = ? OR user_email = ?", login, login)
	if err != nil {
		return user, err
	}
	for rows.Next() {
		var created, lastOnline string
		rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
		date, _ := time.Parse(timeLayout, created)
		user.Created = date
		date, _ = time.Parse(timeLayout, lastOnline)
		user.LastOnline = date
	}
	return user, nil
}
