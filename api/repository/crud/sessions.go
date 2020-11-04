package crud

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
)

func (UserRepoCRUD) ValidateSession(sessionID string) (user models.UserCtx, status int, err error) {
	if err = repository.DB.QueryRow(
		`SELECT id,
			role,
			display_name
		FROM users
		WHERE session_id = ?`, sessionID,
	).Scan(
		&user.ID, &user.Role, &user.DisplayName,
	); err != nil {
		if err != sql.ErrNoRows {
			return user, http.StatusInternalServerError, err
		}
		return models.UserCtx{ID: -1, Role: -1, DisplayName: ""}, http.StatusUnauthorized, errors.New("user not authorized")
	}
	return user, http.StatusOK, nil
}

func (UserRepoCRUD) UpdateSession(id int64, newSession string) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`UPDATE users
		SET session_id = ?
		WHERE id = ?`,
		newSession, id,
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
