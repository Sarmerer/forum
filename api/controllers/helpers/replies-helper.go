package helpers

import (
	"database/sql"
	models "forum/api/models/post"
	"forum/database"
)

func NewReplyModel() (pm *models.PostReplyModel, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	if pm, err = models.NewPostReplyModel(db); err != nil {
		return
	}
	return
}

// func ReplyExists(rid int) (int, error) {
// 	var (
// 		rm     *models.PostReplyModel
// 		status int
// 		err    error
// 	)
// 	if rm, err = NewReplyModel(); err != nil {
// 		return http.StatusInternalServerError, err
// 	}
// 	if _, status, err = rm.FindByID(uid); err != nil {
// 		return status, err
// 	}
// 	return http.StatusOK, nil
// }
