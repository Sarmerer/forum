package models

import (
	"database/sql"
	"errors"
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

func (um *UserModel) GetRole(id uint64) (role int, err error) {
	err = um.DB.QueryRow("SELECT user_role FROM users WHERE user_id = ?", id).Scan(&role)
	if err == sql.ErrNoRows {
		return 0, errors.New("counld not find user with such ID")
	} else if err != nil {
		return 0, err
	}
	return role, nil
}
