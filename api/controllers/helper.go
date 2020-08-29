package controllers

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

func PostExists(pid uint64) (int, error) {
	var (
		db     *sql.DB
		pm     repository.PostRepo
		status int
		err    error
	)
	if db, err = database.Connect(); err != nil {
		return http.StatusInternalServerError, err
	}
	if pm, err = crud.NewPostModel(db); err != nil {
		return http.StatusInternalServerError, err
	}
	if _, status, err = pm.FindByID(pid); err != nil {
		return status, err
	}
	return http.StatusOK, nil
}
