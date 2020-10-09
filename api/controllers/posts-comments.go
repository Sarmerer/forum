package controllers

import (
	"encoding/json"
	"forum/api/config"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/api/utils"
	"net/http"
	"time"
)

func GetComments(w http.ResponseWriter, r *http.Request) {
	var (
		repo     repository.CommentRepo = crud.NewCommentRepoCRUD()
		comments []models.Comment
		pid      int64
		status   int
		err      error
	)
	if pid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if _, status, err = crud.NewPostRepoCRUD().FindByID(pid); err != nil {
		response.Error(w, status, err)
	}
	if comments, err = repo.FindAll(pid); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, nil, comments)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.CommentRepo = crud.NewCommentRepoCRUD()
		uid    int64                  = r.Context().Value("uid").(int64)
		author *models.User
		status int
		err    error
	)
	input := struct {
		PID     int64  `json:"pid"`
		Content string `json:"content"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if _, status, err = crud.NewPostRepoCRUD().FindByID(input.PID); err != nil {
		response.Error(w, status, err)
		return
	}
	if author, status, err = crud.NewUserRepoCRUD().FindByID(uid); err != nil {
		response.Error(w, status, err)
		return
	}
	reply := &models.Comment{
		Content:    input.Content,
		Created:    time.Now().Format(config.TimeLayout),
		Post:       input.PID,
		AuthorID:   uid,
		AuthorName: author.DisplayName,
	}
	if err = repo.Create(reply); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been added", nil)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	var (
		repo         repository.CommentRepo = crud.NewCommentRepoCRUD()
		updatedReply *models.Comment
		status       int
		err          error
	)
	input := struct {
		RID     int64  `json:"id"`
		Content string `json:"content"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if updatedReply, status, err = repo.FindByID(input.RID); err != nil {
		response.Error(w, status, err)
		return
	}

	updatedReply.Content = input.Content
	if err = repo.Update(updatedReply); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been updated", nil)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.CommentRepo = crud.NewCommentRepoCRUD()
		rid    int64
		status int
		err    error
	)
	if rid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if _, status, err = repo.FindByID(rid); err != nil {
		response.Error(w, status, err)
		return
	}
	if err = repo.Delete(rid); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been deleted", nil)
}

func DeleteCommentsGroup(pid int64) error {
	var (
		repo repository.CommentRepo = crud.NewCommentRepoCRUD()
		err  error
	)
	if err = repo.DeleteGroup(pid); err != nil {
		return err
	}
	return nil
}

func CountComments(pid int64) (comments string, err error) {
	var repo = crud.NewCommentRepoCRUD()
	if comments, err = repo.Count(pid); err != nil {
		return "0", err
	}
	return comments, nil
}
