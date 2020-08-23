package models

import (
	"database/sql"
	"errors"
	"forum/api/entities"
	"forum/config"
	"forum/database"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//UserModel helps performing CRUD operations
type UserModel struct {
	DB *sql.DB
}

//NewUserModel creates an instance of UserModel
func NewUserModel() (*sql.DB, *UserModel, error) {
	db, dbErr := database.Connect()
	if dbErr != nil {
		return nil, nil, dbErr
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (user_id INTEGER PRIMARY KEY, user_name TEXT, user_password	BLOB, user_email TEXT, user_nickname	TEXT, user_created	TEXT, user_last_online	TEXT, user_session_id TEXT, user_role INTEGER)")
	if err != nil {
		return nil, nil, err
	}
	statement.Exec()
	return db, &UserModel{db}, nil
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
		date, dateErr := time.Parse(config.TimeLayout, created)
		if dateErr != nil {
			return users, dateErr
		}
		user.Created = date
		date, dateErr = time.Parse(config.TimeLayout, lastOnline)
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
func (um *UserModel) Create(user *entities.User) error {
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
func (um *UserModel) Delete(id int64) error {
	res, err := um.DB.Exec("DELETE FROM users WHERE user_id = ?", id)
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
	return errors.New("unable to delete account")
}

//Update updates existing user in the database
func (um *UserModel) Update(user *entities.User) error {
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

//UpdateRole updates user role in the database
func (um *UserModel) UpdateRole(userID, role int) error {
	statement, err := um.DB.Prepare("UPDATE users SET user_role = ? WHERE user_id = ?")
	if err != nil {
		return err
	}
	res, err := statement.Exec(role, userID)
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
	return errors.New("failed to update role")
}

func (um *UserModel) GetRole(id int64) (role int, err error) {
	err = um.DB.QueryRow("SELECT user_role FROM users WHERE user_id = ?", id).Scan(&role)
	if err == sql.ErrNoRows {
		return 0, errors.New("counld not find user with such ID")
	} else if err != nil {
		return 0, err
	}
	return role, nil
}

//FindByNameOrEmail finds a user by name or email in the database
func (um *UserModel) FindByNameOrEmail(login string) (entities.User, error) {
	var user entities.User
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
