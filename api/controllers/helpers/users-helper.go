package helpers

import (
	"database/sql"
	models "forum/api/models/user"
	"forum/database"
	"net/http"
)

func NewUserModel() (pm *models.UserModel, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	if pm, err = models.NewUserModel(db); err != nil {
		return
	}
	return
}

func UserExists(uid uint64) (int, error) {
	var (
		um     *models.UserModel
		status int
		err    error
	)
	if um, err = NewUserModel(); err != nil {
		return http.StatusInternalServerError, err
	}
	if _, status, err = um.FindByID(uid); err != nil {
		return status, err
	}
	return http.StatusOK, nil
}
