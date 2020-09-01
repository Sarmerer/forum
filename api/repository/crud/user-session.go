package crud

import (
	"database/sql"
	"errors"
)

func (um *UserModel) ValidateSession(session string) (uid uint64, err error) {
	err = um.DB.QueryRow("SELECT id FROM users WHERE session_id = ?", session).Scan(&uid)
	if err == sql.ErrNoRows {
		err = errors.New("session not found")
		return
	} else if err != nil {
		return
	}
	return
}

func (um *UserModel) UpdateSession(id uint64, newSession string) error {
	statement, err := um.DB.Prepare("UPDATE users SET session_id = ? WHERE id = ?")
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
