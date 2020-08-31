package controllers

import (
	"database/sql"
	"encoding/json"
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

	if prm, err = helpers.PrepareReplyRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if replies, err = prm.FindAll(pid); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, nil, replies)
}

func GetReply(w http.ResponseWriter, r *http.Request) {
	var (
		rid    uint64
		prm    repository.ReplyRepo
		reply  *models.PostReply
		status int
		err    error
	)
	if rid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if prm, err = helpers.PrepareReplyRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if reply, status, err = prm.FindByID(rid); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, nil, reply)
}

func CreateReply(w http.ResponseWriter, r *http.Request) {
	var (
		pid    uint64
		prm    repository.ReplyRepo
		author uint64
		status int
		err    error
	)
	input := struct {
		Content string `json:"content"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if pid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if status, err = postExists(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	if prm, err = helpers.PrepareReplyRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	author = r.Context().Value("uid").(uint64)
	reply := &models.PostReply{
		Content: input.Content,
		Date:    time.Now(),
		Post:    pid,
		By:      author,
	}
	if err = prm.Create(reply); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been added", nil)
}

func UpdateReply(w http.ResponseWriter, r *http.Request) {
	var (
		rid          uint64
		prm          repository.ReplyRepo
		updatedReply *models.PostReply
		status       int
		err          error
	)
	input := struct {
		Content string `json:"content"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if rid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if prm, err = helpers.PrepareReplyRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if updatedReply, status, err = prm.FindByID(rid); err != nil {
		response.Error(w, status, err)
		return
	}
	if input.Content != "" {
		updatedReply.Content = input.Content
	}
	if err = prm.Update(updatedReply); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been updated", nil)
}

func DeleteReply(w http.ResponseWriter, r *http.Request) {
	var (
		rid    uint64
		prm    repository.ReplyRepo
		status int
		err    error
	)
	if rid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if prm, err = helpers.PrepareReplyRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if _, status, err = prm.FindByID(rid); err != nil {
		response.Error(w, status, err)
		return
	}
	if err = prm.Delete(rid); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been deleted", nil)
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
	pm = crud.NewPostModel(db)
	if _, status, err = pm.FindByID(pid); err != nil {
		return status, err
	}
	return http.StatusOK, nil
}
