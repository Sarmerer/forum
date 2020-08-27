package controllers

import (
	"database/sql"
	"errors"
	models "forum/api/models/post"
	"forum/api/response"
	"forum/database"
	"net/http"
	"strconv"
	"time"
)

func GetReplies(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sql.DB
	var replies []models.PostReply
	pid, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("invalid ID parameter"))
		return
	}
	db, err = database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}
	pr, err := models.NewPostReplyModel(db)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}
	replies, err = pr.FindAll(pid)
	if err != nil {
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
		}
	}
	response.Success(w, nil, replies)
}

func CreateReply(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sql.DB
	var status int
	pid, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("invalid ID parameter"))
		return
	}
	db, err = database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}
	pr, err := models.NewPostReplyModel(db)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}
	reply := &models.PostReply{
		Content: "test content",
		Date:    time.Now(),
		Post:    pid,
		By:      0,
	}
	status, err = pr.Create(reply)
	if err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, "reply has been added", nil)
}

func UpdateReply(w http.ResponseWriter, r *http.Request) {

}

func DeleteReply(w http.ResponseWriter, r *http.Request) {

}
