package crud

import (
	"database/sql"
	"errors"
	"forum/api/models"
	"forum/config"
	"net/http"
	"time"
)

//UserModel helps performing CRUD operations
type UserModel struct {
	DB *sql.DB
}

//NewUserModel creates an instance of UserModel
func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{db}
}

//FindAll returns all users in the database
func (um *UserModel) FindAll() ([]models.User, error) {
	var err error
	var users []models.User
	var rows *sql.Rows
	rows, err = um.DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user models.User
		var created, lastOnline string
		rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &created, &lastOnline, &user.SessionID, &user.Role)
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
func (um *UserModel) FindByID(uid uint64) (*models.User, int, error) {
	var user models.User
	var created, lastOnline string
	row := um.DB.QueryRow("SELECT * FROM users WHERE id = ?", uid)
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &created, &lastOnline, &user.SessionID, &user.Role)
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
func (um *UserModel) Create(user *models.User) (int, error) {
	statement, err := um.DB.Prepare("INSERT INTO users (name, password, email, created, last_online, session_id, role) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	nameErr := um.DB.QueryRow("SELECT name FROM users WHERE name = ?", user.Name).Scan(&user.Name)
	emailErr := um.DB.QueryRow("SELECT email FROM users WHERE email = ?", user.Email).Scan(&user.Email)
	if nameErr != sql.ErrNoRows && nameErr != nil || emailErr != sql.ErrNoRows && emailErr != nil {
		return http.StatusInternalServerError, errors.New("unable to create new user account")
	} else if nameErr == sql.ErrNoRows && emailErr == sql.ErrNoRows {
		res, err := statement.Exec(
			user.Name, user.Password, user.Email, time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout), user.SessionID, user.Role,
		)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if rowsAffected > 0 {
			return http.StatusOK, nil
		}
	} else if nameErr != sql.ErrNoRows && emailErr == sql.ErrNoRows {
		return http.StatusConflict, errors.New("login already taken")
	} else if emailErr != sql.ErrNoRows && nameErr == sql.ErrNoRows {
		return http.StatusConflict, errors.New("email already taken")
	} else if nameErr == nil && emailErr == nil {
		return http.StatusConflict, errors.New("login and email are already taken")
	}
	return http.StatusBadRequest, errors.New("unable to create new user account")
}

//Update updates existing user in the database
func (um *UserModel) Update(user *models.User) (int, error) {
	statement, err := um.DB.Prepare("UPDATE users SET name = ?, email = ?, last_online = ? WHERE id = ?")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	res, err := statement.Exec(user.Name, user.Email, user.LastOnline.Format(config.TimeLayout), user.ID)
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
	var result sql.Result
	var rowsAffected int64
	result, err = um.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return http.StatusBadRequest, errors.New("user not found")
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); rowsAffected > 0 && err == nil {
		return http.StatusOK, nil
	}
	return http.StatusInternalServerError, err
}

//FindByNameOrEmail finds a user by name or email in the database
func (um *UserModel) FindByNameOrEmail(login string) (*models.User, int, error) {
	var user models.User
	rows, err := um.DB.Query("SELECT * FROM users WHERE name = ? OR email = ?", login, login)
	if err == sql.ErrNoRows {
		return nil, http.StatusBadRequest, errors.New("wrong login or password")
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	for rows.Next() {
		var created, lastOnline string
		rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &created, &lastOnline, &user.SessionID, &user.Role)
		if user.Created, err = time.Parse(config.TimeLayout, created); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if user.LastOnline, err = time.Parse(config.TimeLayout, lastOnline); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}
	return &user, http.StatusOK, nil
}
