package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"
)

func GetComments(w http.ResponseWriter, r *http.Request) {
	var (
		repo     repository.CommentRepo = crud.NewCommentRepoCRUD()
		userCtx  models.UserCtx         = utils.GetUserFromCtx(r)
		comments []models.Comment
		pid      int64
		status   int
		err      error
	)
	if pid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if _, status, err = crud.NewPostRepoCRUD().FindByID(pid, -1); err != nil {
		response.Error(w, status, err)
		return
	}
	if comments, err = repo.FindAll(pid, userCtx.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, nil, comments)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.CommentRepo = crud.NewCommentRepoCRUD()
		userCtx models.UserCtx         = utils.GetUserFromCtx(r)
		status  int
		err     error
	)
	input := struct {
		PID     int64  `json:"pid"`
		Content string `json:"content"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if _, status, err = crud.NewPostRepoCRUD().FindByID(input.PID, -1); err != nil {
		response.Error(w, status, err)
		return
	}
	comment := &models.Comment{
		Content:    input.Content,
		Created:    time.Now().Format(config.TimeLayout),
		PostID:     input.PID,
		AuthorID:   userCtx.ID,
		AuthorName: userCtx.DisplayName,
	}
	if err = repo.Create(comment); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been added", nil)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.CommentRepo = crud.NewCommentRepoCRUD()
		comment *models.Comment
		status  int
		err     error
	)
	input := struct {
		RID     int64  `json:"id"`
		Content string `json:"content"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if comment, status, err = repo.FindByID(input.RID); err != nil {
		response.Error(w, status, err)
		return
	}

	if !requestorIsEntityOwner(utils.GetUserFromCtx(r), comment.AuthorID) {
		response.Error(w, http.StatusForbidden, errors.New("this comment doesn't belong to you"))
		return
	}

	comment.Content = input.Content
	if err = repo.Update(comment); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been updated", nil)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.CommentRepo = crud.NewCommentRepoCRUD()
		cid     int64
		comment *models.Comment
		status  int
		err     error
	)
	if cid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if comment, status, err = repo.FindByID(cid); err != nil {
		response.Error(w, status, err)
		return
	}

	if !requestorIsEntityOwner(utils.GetUserFromCtx(r), comment.AuthorID) {
		response.Error(w, http.StatusForbidden, errors.New("this comment doesn't belong to you"))
		return
	}

	if err = repo.Delete(cid); err != nil {
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
