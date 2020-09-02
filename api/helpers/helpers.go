package helpers

import (
	"database/sql"
	"errors"
	"forum/api/repository"
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

func PrepareUserRepo() (um repository.UserRepo, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	um = crud.NewUserModel(db)
	return
}

func PreparePostRepo() (pm repository.PostRepo, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	pm = crud.NewPostModel(db)
	return
}

func PrepareReplyRepo() (prm repository.ReplyRepo, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	prm = crud.NewReplyRepoCRUD(db)
	return
}

func PrepareCategoriesRepo() (cm repository.CategoryRepo, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	cm = crud.NewCategoryModel(db)
	return
}
