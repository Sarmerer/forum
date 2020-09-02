package crud

import (
	"database/sql"
	"errors"
	"forum/database"
	"net/http"
)

func (um *UserRepoCRUD) ValidateSession(session string) (uid uint64, status int, err error) {
	if err = database.DB.QueryRow("SELECT id FROM users WHERE session_id = ?", session).Scan(&uid); err != nil {
		if err != sql.ErrNoRows {
			return uid, http.StatusInternalServerError, err
		}
		return 0, http.StatusUnauthorized, errors.New("user not authorized")
	}
	return uid, http.StatusOK, nil
}

func (um *UserRepoCRUD) UpdateSession(id uint64, newSession string) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = database.DB.Exec(
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
