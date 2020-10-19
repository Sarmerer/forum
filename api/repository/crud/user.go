package crud

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
)

//UserRepoCRUD helps performing CRUD operations
type UserRepoCRUD struct{}

//NewUserRepoCRUD creates an instance of UserModel
func NewUserRepoCRUD() UserRepoCRUD {
	return UserRepoCRUD{}
}

//FindAll returns all users in the database
func (UserRepoCRUD) FindAll() ([]models.User, error) {
	var (
		rows  *sql.Rows
		users []models.User
		err   error
	)
	if rows, err = repository.DB.Query(
		"SELECT * FROM users",
	); err != nil {
		return nil, err
	}
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Login, &u.Password, &u.Email, &u.Avatar, &u.DisplayName, &u.Created, &u.LastOnline, &u.SessionID, &u.Role)
		users = append(users, u)
	}
	return users, nil
}

//FindByID returns a specific user from the database
//FIXME don't scan for sensetive data, like password and session id
func (UserRepoCRUD) FindByID(uid int64) (*models.User, int, error) {
	var (
		u   models.User
		err error
	)
	if err = repository.DB.QueryRow(
		"SELECT * FROM users WHERE id = ?", uid,
	).Scan(
		&u.ID, &u.Login, &u.Password, &u.Email, &u.Avatar, &u.DisplayName, &u.Created, &u.LastOnline, &u.SessionID, &u.Role,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, errors.New("user not found")
	}
	return &u, http.StatusOK, nil
}

//Create adds a new user to the database
func (UserRepoCRUD) Create(u *models.User) (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	name := repository.DB.QueryRow("SELECT login FROM users WHERE login = ?", u.Login).Scan(&u.Login)
	email := repository.DB.QueryRow("SELECT email FROM users WHERE email = ?", u.Email).Scan(&u.Email)
	if name == nil && email != nil {
		return http.StatusConflict, errors.New("name is not unique")
	} else if email == nil && name != nil {
		return http.StatusConflict, errors.New("email is not unique")
	} else if name == nil && email == nil {
		return http.StatusConflict, errors.New("name and email are not unique")
	}

	if result, err = repository.DB.Exec(
		"INSERT INTO users (login, password, email, avatar, display_name, created, last_online, session_id, role) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		u.Login, u.Password, u.Email, u.Avatar, u.DisplayName, time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout), u.SessionID, u.Role,
	); err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not create the user")
}

//Update updates existing user in the database
//TODO decide what we will let users to update
func (UserRepoCRUD) Update(u *models.User) (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"UPDATE users SET display_name = ? WHERE id = ?",
		u.DisplayName, u.ID,
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not update the user")
}

//Delete deletes user from the database
func (UserRepoCRUD) Delete(id int64) (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"DELETE FROM users WHERE id = ?", id,
	); err != nil {
		if err != sql.ErrNoRows {
			return http.StatusInternalServerError, err
		}
		return http.StatusNotFound, errors.New("user not found")
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not delete the user")
}

//FindByNameOrEmail finds a user by name or email in the database
func (UserRepoCRUD) FindByNameOrEmail(login string) (*models.User, int, error) {
	var (
		u   models.User
		err error
	)
	if err = repository.DB.QueryRow(
		"SELECT * FROM users WHERE login = ? OR email = ?", login, login,
	).Scan(
		&u.ID, &u.Login, &u.Password, &u.Email, &u.Avatar, &u.DisplayName, &u.Created, &u.LastOnline, &u.SessionID, &u.Role,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusBadRequest, errors.New("wrong login or password")
	}
	return &u, http.StatusOK, nil
}
