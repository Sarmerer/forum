package crud

import (
	"database/sql"
	"errors"
	"forum/api/repository"
	"net/http"
)

func (UserRepoCRUD) ValidateSession(session string) (uid int64, role, status int, err error) {
	if err = repository.DB.QueryRow("SELECT id, role FROM users WHERE session_id = ?", session).Scan(&uid, &role); err != nil {
		if err != sql.ErrNoRows {
			return uid, role, http.StatusInternalServerError, err
		}
		return -1, -1, http.StatusUnauthorized, errors.New("user not authorized")
	}
	return uid, role, http.StatusOK, nil
}

func (UserRepoCRUD) UpdateSession(id int64, newSession string) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"UPDATE users SET session_id = ? WHERE id = ?",
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
