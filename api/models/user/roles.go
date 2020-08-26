package models

import (
	"database/sql"
	"errors"
	"net/http"
)

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

func (um *UserModel) GetRole(id uint64) (int, int, error) {
	var role int
	err := um.DB.QueryRow("SELECT user_role FROM users WHERE user_id = ?", id).Scan(&role)
	if err == sql.ErrNoRows {
		return 0, http.StatusBadRequest, errors.New("user not found")
	} else if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	return role, http.StatusOK, nil
}
