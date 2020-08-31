package helpers

import (
	"database/sql"
	"errors"
	"forum/api/repository/crud"
	"forum/database"
	"net/http"
	"strconv"
)

func ParseID(r *http.Request) (res uint64, err error) {
	if res, err = strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64); err != nil {
		return 0, errors.New("invalid id")
	}
	return res, nil
}

func PrepareUserRepo() (um *crud.UserModel, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	um = crud.NewUserModel(db)
	return
}

func PreparePostRepo() (pm *crud.PostModel, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	pm = crud.NewPostModel(db)
	return
}

func PrepareReplyRepo() (prm *crud.PostReplyModel, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	prm = crud.NewPostReplyModel(db)
	return
}
