package models

import (
	"errors"
)

func (um *UserModel) ValidateSession() error {
	return errors.New("f")
}

func (um *UserModel) UpdateSession(id uint64, newSession string) error {
	statement, err := um.DB.Prepare("UPDATE users SET user_session_id = ? WHERE user_id = ?")
	if err != nil {
		return err
	}
	res, err := statement.Exec(newSession, id)
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
	return errors.New("failed to update the session")
}
