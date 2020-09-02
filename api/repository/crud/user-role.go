package crud

import (
	"database/sql"
	"errors"
	"forum/database"
	"net/http"
)

func (repo *UserRepoCRUD) GetRole(id uint64) (int, int, error) {
	var (
		role int
		err  error
	)
	if err = database.DB.QueryRow(
		"SELECT role FROM users WHERE id = ?", id,
	).Scan(&role); err != nil {
		if err != sql.ErrNoRows {
			return 0, http.StatusInternalServerError, err
		}
		return 0, http.StatusNoContent, errors.New("user not found")
	}
	return role, http.StatusOK, nil
}

//UpdateRole updates user role in the database
func (um *UserRepoCRUD) UpdateRole(uid uint64, role int) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = database.DB.Exec(
		"UPDATE users SET role = ? WHERE id = ?",
		role, uid,
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
