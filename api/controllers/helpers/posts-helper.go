package helpers

import (
	"database/sql"
	models "forum/api/models/post"
	"forum/database"
	"net/http"
)

func NewPostModel() (pm *models.PostModel, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	if pm, err = models.NewPostModel(db); err != nil {
		return
	}
	return
}

func PostExists(pid uint64) (int, error) {
	var (
		pm     *models.PostModel
		status int
		err    error
	)
	if pm, err = NewPostModel(); err != nil {
		return http.StatusInternalServerError, err
	}
	if _, status, err = pm.FindByID(pid); err != nil {
		return status, err
	}
	return http.StatusOK, nil
}
