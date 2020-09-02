package crud

import (
	"database/sql"
	"errors"
	"forum/api/models"
	"forum/config"
	"forum/database"
	"net/http"
	"time"
)

//UserRepoCRUD helps performing CRUD operations
type UserRepoCRUD struct{}

//NewUserRepoCRUD creates an instance of UserModel
func NewUserRepoCRUD() *UserRepoCRUD {
	return &UserRepoCRUD{}
}

//FindAll returns all users in the database
func (repo *UserRepoCRUD) FindAll() ([]models.User, error) {
	var (
		rows  *sql.Rows
		users []models.User
		err   error
	)
	if rows, err = database.DB.Query(
		"SELECT * FROM users",
	); err != nil {
		return nil, err
	}
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name, &u.Password, &u.Email, &u.Created, &u.LastOnline, &u.SessionID, &u.Role)
		users = append(users, u)
	}
	return users, nil
}

//FindByID returns a specific user from the database
func (repo *UserRepoCRUD) FindByID(uid uint64) (*models.User, int, error) {
	var (
		u   models.User
		err error
	)
	if err = database.DB.QueryRow(
		"SELECT * FROM users WHERE id = ?", uid,
	).Scan(
		&u.ID, &u.Name, &u.Password, &u.Email, &u.Created, &u.LastOnline, &u.SessionID, &u.Role,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNoContent, errors.New("user not found")
	}
	return &u, http.StatusOK, nil
}

//Create adds a new user to the database
func (repo *UserRepoCRUD) Create(u *models.User) (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	name := database.DB.QueryRow("SELECT name FROM users WHERE name = ?", u.Name).Scan(&u.Name)
	email := database.DB.QueryRow("SELECT email FROM users WHERE email = ?", u.Email).Scan(&u.Email)
	if name == nil && email != nil {
		return http.StatusConflict, errors.New("name is not unique")
	} else if email == nil && name != nil {
		return http.StatusConflict, errors.New("email is not unique")
	} else if name == nil && email == nil {
		return http.StatusConflict, errors.New("name and email are not unique")
	}

	if result, err = database.DB.Exec(
		"INSERT INTO users (name, password, email, created, last_online, session_id, role) VALUES (?, ?, ?, ?, ?, ?, ?)",
		u.Name, u.Password, u.Email, time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout), u.SessionID, u.Role,
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
//TODO check if new naame is unique
func (repo *UserRepoCRUD) Update(u *models.User) (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = database.DB.Exec(
		"UPDATE users SET name = ?, email = ?, last_online = ? WHERE id = ?",
		u.Name, u.Email, u.LastOnline, u.ID,
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
func (repo *UserRepoCRUD) Delete(id uint64) (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = database.DB.Exec(
		"DELETE FROM users WHERE id = ?", id,
	); err != nil {
		if err != sql.ErrNoRows {
			return http.StatusInternalServerError, err
		}
		return http.StatusNoContent, errors.New("user not found")
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
func (repo *UserRepoCRUD) FindByNameOrEmail(login string) (*models.User, int, error) {
	var (
		u   models.User
		err error
	)
	if err = database.DB.QueryRow(
		"SELECT * FROM users WHERE name = ? OR email = ?", login, login,
	).Scan(
		&u.ID, &u.Name, &u.Password, &u.Email, &u.Created, &u.LastOnline, &u.SessionID, &u.Role,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusBadRequest, errors.New("wrong login or password")
	}
	return &u, http.StatusOK, nil
}
