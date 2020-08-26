package models

import (
	"database/sql"
	"errors"
	"forum/api/utils/channel"
	"forum/config"
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
func (um *UserModel) FindAll() (users []User, err error) {
	var rows *sql.Rows
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rows, err = um.DB.Query("SELECT * FROM users")
		if err != nil {
			ch <- false
			return
		}
		for rows.Next() {
			var date time.Time
			var user User
			var created, lastOnline string
			rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
			date, err := time.Parse(config.TimeLayout, created)
			if err != nil {
				ch <- false
				return
			}
			user.Created = date
			date, err = time.Parse(config.TimeLayout, lastOnline)
			if err != nil {
				ch <- false
				return
			}
			user.LastOnline = date
			users = append(users, user)
		}
		ch <- true
	}(done)
	if channel.OK(done) {
		return users, nil
	}
	return nil, err
}

//FindByID returns a specific user from the database
func (um *UserModel) FindByID(id uint64) (User, error) {
	var user User
	rows, err := um.DB.Query("SELECT * FROM users WHERE user_id = ?", id)
	if err != nil {
		return user, err
	}
	for rows.Next() {
		var created, lastOnline string
		rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
		date, err := time.Parse(config.TimeLayout, created)
		if err != nil {
			return user, err
		}
		user.Created = date
		date, err = time.Parse(config.TimeLayout, lastOnline)
		if err != nil {
			return user, err
		}
		user.LastOnline = date
	}
	return user, nil
}

//Create adds a new user to the database
func (um *UserModel) Create(user *User) error {
	statement, err := um.DB.Prepare("INSERT INTO users (user_name, user_password, user_email, user_nickname, user_created, user_last_online, user_session_id, user_role) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("unable to create new user account")
	}
	nameErr := um.DB.QueryRow("SELECT user_name FROM users WHERE user_name = ?", user.Name).Scan(&user.Name)
	emailErr := um.DB.QueryRow("SELECT user_email FROM users WHERE user_email = ?", user.Email).Scan(&user.Email)
	if nameErr != sql.ErrNoRows && nameErr != nil || emailErr != sql.ErrNoRows && emailErr != nil {
		return errors.New("unable to create new user account")
	} else if nameErr == sql.ErrNoRows && emailErr == sql.ErrNoRows {
		res, err := statement.Exec(
			user.Name, user.Password, user.Email, user.Nickname, time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout), user.SessionID, user.Role,
		)
		if err != nil {
			return errors.New("unable to create new user account")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return errors.New("unable to create new user account")
		}
		if rowsAffected > 0 {
			return nil
		}
	} else if nameErr != sql.ErrNoRows && emailErr == sql.ErrNoRows {
		return errors.New("name not unique")
	} else if emailErr != sql.ErrNoRows && nameErr == sql.ErrNoRows {
		return errors.New("email not unique")
	} else if nameErr == nil && emailErr == nil {
		return errors.New("both not unique")
	}
	return errors.New("unable to create new user account")
}

//Delete deletes user from the database
func (um *UserModel) Delete(id uint64) error {
	var err error
	var response sql.Result
	var rowsAffected int64
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		response, err = um.DB.Exec("DELETE FROM users WHERE user_id = ?", id)
		if err != nil {
			ch <- false
			return
		}
		if rowsAffected, err = response.RowsAffected(); rowsAffected > 0 && err == nil {
			ch <- true
		} else {
			ch <- false
			return
		}
	}(done)
	if channel.OK(done) {
		return nil
	}
	return err
}

//Update updates existing user in the database
func (um *UserModel) Update(user *User) error {
	statement, err := um.DB.Prepare("UPDATE users SET user_name = ?, user_email = ?, user_nickname = ?, user_last_online = ? WHERE user_id = ?")
	if err != nil {
		return err
	}
	res, err := statement.Exec(user.Name, user.Email, user.Nickname, user.LastOnline.Format(config.TimeLayout), user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("failed to update the user")
}

//FindByNameOrEmail finds a user by name or email in the database
func (um *UserModel) FindByNameOrEmail(login string) (User, error) {
	var user User
	rows, err := um.DB.Query("SELECT * FROM users WHERE user_name = ? OR user_email = ?", login, login)
	if err != nil {
		return user, err
	}
	for rows.Next() {
		var created, lastOnline string
		rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
		date, _ := time.Parse(config.TimeLayout, created)
		user.Created = date
		date, _ = time.Parse(config.TimeLayout, lastOnline)
		user.LastOnline = date
	}
	return user, nil
}
