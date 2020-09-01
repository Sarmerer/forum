package crud

import (
	"database/sql"
	"errors"
	"net/http"
)

func (um *UserModel) GetRole(id uint64) (int, int, error) {
	var role int
	err := um.DB.QueryRow("SELECT role FROM users WHERE id = ?", id).Scan(&role)
	if err == sql.ErrNoRows {
		return 0, http.StatusBadRequest, errors.New("user not found")
	} else if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	return role, http.StatusOK, nil
}

//UpdateRole updates user role in the database
func (um *UserModel) UpdateRole(uid uint64, role int) error {
	statement, err := um.DB.Prepare("UPDATE users SET role = ? WHERE id = ?")
	if err != nil {
		return err
	}
	res, err := statement.Exec(role, uid)
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
