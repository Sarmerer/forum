package controllers

import (
	"database/sql"
	"forum/api/controllers/helpers"
	models "forum/api/models/post"
	"forum/api/response"
	"forum/database"
	"net/http"
	"time"
)

func GetReplies(w http.ResponseWriter, r *http.Request) {
	var (
		pid     uint64
		pr      *models.PostReplyModel
		replies []models.PostReply
		status  int
		err     error
	)
	if pid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if status, err = helpers.PostExists(pid); err != nil {
		response.Error(w, status, err)
		return
	}

	if pr, err = helpers.NewReplyModel(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if replies, err = pr.FindAll(pid); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, nil, replies)
}

func CreateReply(w http.ResponseWriter, r *http.Request) {
	var (
		pid    uint64
		db     *sql.DB
		pr     *models.PostReplyModel
		author uint64
		status int
		err    error
	)
	if pid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	db, err = database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}
	if status, err = helpers.PostExists(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	if pr, err = models.NewPostReplyModel(db); err != nil {
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
	if status, err = pr.Create(reply); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, "reply has been added", nil)
}

func UpdateReply(w http.ResponseWriter, r *http.Request) {

}

func DeleteReply(w http.ResponseWriter, r *http.Request) {

}
