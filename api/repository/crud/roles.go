package crud

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/sarmerer/forum/api/repository"
)

func (UserRepoCRUD) GetRole(userID int64) (int, int, error) {
	var (
		role int
		err  error
	)
	if err = repository.DB.QueryRow(
		`SELECT role
		FROM users
		WHERE _id = ?`, userID,
	).Scan(&role); err != nil {
		if err != sql.ErrNoRows {
			return 0, http.StatusInternalServerError, err
		}
		return 0, http.StatusNotFound, errors.New("user not found")
	}
	return role, http.StatusOK, nil
}

//UpdateRole updates user role in the database
func (UserRepoCRUD) UpdateRole(userID int64, role int) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`UPDATE users
		SET role = $1
		WHERE _id = $2`,
		role, userID,
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
