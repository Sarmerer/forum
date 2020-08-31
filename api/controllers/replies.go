package controllers

import (
	"database/sql"
	"forum/api/helpers"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/database"
	"net/http"
	"time"
)

func GetReplies(w http.ResponseWriter, r *http.Request) {
	var (
		pid     uint64
		prm     repository.ReplyRepo
		replies []models.PostReply
		status  int
		err     error
	)
	if pid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if status, err = postExists(pid); err != nil {
		response.Error(w, status, err)
		return
	}

	if prm, err = newPRM(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if replies, err = prm.FindAll(pid); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, nil, replies)
}

func CreateReply(w http.ResponseWriter, r *http.Request) {
	var (
		pid    uint64
		prm    repository.ReplyRepo
		author uint64
		status int
		err    error
	)
	if pid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if status, err = postExists(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	if prm, err = newPRM(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	author = r.Context().Value("uid").(uint64)
	reply := &models.PostReply{
		Content: "test content",
		Date:    time.Now(),
		Post:    pid,
		By:      author,
	}
	if status, err = prm.Create(reply); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, "reply has been added", nil)
}

func UpdateReply(w http.ResponseWriter, r *http.Request) {

}

func DeleteReply(w http.ResponseWriter, r *http.Request) {

}

func newPRM() (prm *crud.PostReplyModel, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	if prm, err = crud.NewPostReplyModel(db); err != nil {
		return
	}
	return
}

func postExists(pid uint64) (int, error) {
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
