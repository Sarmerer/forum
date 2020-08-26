package models

import (
	"database/sql"
	"errors"
	"forum/config"
	"net/http"
	"time"
)

type User struct {
	ID         uint64
	Name       string
	Password   string
	Email      string
	Nickname   string
	Created    time.Time
	LastOnline time.Time
	SessionID  string
	Role       int
}

//UserModel helps performing CRUD operations
type UserModel struct {
	DB *sql.DB
}

//NewUserModel creates an instance of UserModel
func NewUserModel(db *sql.DB) (*UserModel, error) {
	return &UserModel{db}, nil
}

//FindAll returns all users in the database
func (um *UserModel) FindAll() ([]User, error) {
	var err error
	var users []User
	var rows *sql.Rows
	rows, err = um.DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		var created, lastOnline string
		rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
		if user.Created, err = time.Parse(config.TimeLayout, created); err != nil {
			return nil, err
		}
		if user.LastOnline, err = time.Parse(config.TimeLayout, lastOnline); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

//FindByID returns a specific user from the database
func (um *UserModel) FindByID(uid uint64) (*User, int, error) {
	var user User
	var created, lastOnline string
	row := um.DB.QueryRow("SELECT * FROM users WHERE user_id = ?", uid)
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
	if err == sql.ErrNoRows {
		return nil, http.StatusBadRequest, errors.New("user not found")
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if user.Created, err = time.Parse(config.TimeLayout, created); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if user.LastOnline, err = time.Parse(config.TimeLayout, lastOnline); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &user, http.StatusOK, nil
}

//Create adds a new user to the database
func (um *UserModel) Create(user *User) (int, error) {
	statement, err := um.DB.Prepare("INSERT INTO users (user_name, user_password, user_email, user_nickname, user_created, user_last_online, user_session_id, user_role) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return http.StatusInternalServerError, err //errors.New("unable to create new user account")
	}
	nameErr := um.DB.QueryRow("SELECT user_name FROM users WHERE user_name = ?", user.Name).Scan(&user.Name)
	emailErr := um.DB.QueryRow("SELECT user_email FROM users WHERE user_email = ?", user.Email).Scan(&user.Email)
	if nameErr != sql.ErrNoRows && nameErr != nil || emailErr != sql.ErrNoRows && emailErr != nil {
		return http.StatusInternalServerError, errors.New("unable to create new user account")
	} else if nameErr == sql.ErrNoRows && emailErr == sql.ErrNoRows {
		res, err := statement.Exec(
			user.Name, user.Password, user.Email, user.Nickname, time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout), user.SessionID, user.Role,
		)
		if err != nil {
			return http.StatusInternalServerError, err //errors.New("unable to create new user account")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return http.StatusInternalServerError, err //errors.New("unable to create new user account")
		}
		if rowsAffected > 0 {
			return http.StatusOK, nil
		}
	} else if nameErr != sql.ErrNoRows && emailErr == sql.ErrNoRows {
		return http.StatusConflict, errors.New("name not unique")
	} else if emailErr != sql.ErrNoRows && nameErr == sql.ErrNoRows {
		return http.StatusConflict, errors.New("email not unique")
	} else if nameErr == nil && emailErr == nil {
		return http.StatusConflict, errors.New("both not unique")
	}
	return http.StatusBadRequest, errors.New("unable to create new user account")
}

//Update updates existing user in the database
func (um *UserModel) Update(user *User) (int, error) {
	statement, err := um.DB.Prepare("UPDATE users SET user_name = ?, user_email = ?, user_nickname = ?, user_last_online = ? WHERE user_id = ?")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	res, err := statement.Exec(user.Name, user.Email, user.Nickname, user.LastOnline.Format(config.TimeLayout), user.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return http.StatusInternalServerError, err
	} else if rowsAffected > 0 {
		return http.StatusOK, nil
	}
	return http.StatusBadRequest, errors.New("failed to update the user")
}

//Delete deletes user from the database
func (um *UserModel) Delete(id uint64) (int, error) {
	var err error
	var response sql.Result
	var rowsAffected int64
	response, err = um.DB.Exec("DELETE FROM users WHERE user_id = ?", id)
	if err == sql.ErrNoRows {
		return http.StatusBadRequest, errors.New("user not found")
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err = response.RowsAffected(); rowsAffected > 0 && err == nil {
		return http.StatusOK, nil
	} else {
		return http.StatusInternalServerError, err
	}
}

//FindByNameOrEmail finds a user by name or email in the database
func (um *UserModel) FindByNameOrEmail(login string) (*User, int, error) {
	var user User
	rows, err := um.DB.Query("SELECT * FROM users WHERE user_name = ? OR user_email = ?", login, login)
	if err == sql.ErrNoRows {
		return nil, http.StatusBadRequest, errors.New("wrong login or password")
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	for rows.Next() {
		var created, lastOnline string
		rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
		if user.Created, err = time.Parse(config.TimeLayout, created); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if user.LastOnline, err = time.Parse(config.TimeLayout, lastOnline); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}
	return &user, http.StatusOK, nil
}
