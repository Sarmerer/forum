package crud

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
)

func (UserRepoCRUD) GetPassword(userID int64) (string, int, error) {
	var (
		password string
		err      error
	)
	if err = repository.DB.QueryRow(
		`SELECT password
		FROM users
		WHERE _id = ?`, userID,
	).Scan(
		&password,
	); err != nil {
		if err != sql.ErrNoRows {
			return "", http.StatusInternalServerError, err
		}
		return "", http.StatusNotFound, errors.New("user not found")
	}
	return password, http.StatusOK, nil
}

func (UserRepoCRUD) ValidateSession(sessionID string) (user models.UserCtx, status int, err error) {
	if err = repository.DB.QueryRow(
		`SELECT _id,
			role
		FROM users
		WHERE session_id = ? AND session_id <> ''`, sessionID,
	).Scan(
		&user.ID, &user.Role,
	); err != nil {
		if err != sql.ErrNoRows {
			return user, http.StatusInternalServerError, err
		}
		return models.UserCtx{ID: -1, Role: -1}, http.StatusUnauthorized, errors.New("user not authorized")
	}
	return user, http.StatusOK, nil
}

func (UserRepoCRUD) UpdateSession(userID int64, newSession string) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`UPDATE users
		SET session_id = ?
		WHERE _id = ?`,
		newSession, userID,
	); err != nil {
		return err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return nil
}
