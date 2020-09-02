package controllers

import (
	"encoding/json"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/api/utils"
	"forum/config"
	"net/http"
	"time"
)

func GetReplies(pid uint64) ([]models.PostReply, error) {
	var (
		repo    repository.ReplyRepo = crud.NewReplyRepoCRUD()
		replies []models.PostReply
		err     error
	)

	if replies, err = repo.FindAll(pid); err != nil {
		return nil, err
	}
	return replies, nil
}

func CreateReply(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.ReplyRepo = crud.NewReplyRepoCRUD()
		author uint64               = r.Context().Value("uid").(uint64)
		pid    uint64
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
	if pid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if _, status, err = crud.NewPostRepoCRUD().FindByID(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	reply := &models.PostReply{
		Content: input.Content,
		Created: time.Now().Format(config.TimeLayout),
		Post:    pid,
		Author:  author,
	}
	if err = repo.Create(reply); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been added", nil)
}

func UpdateReply(w http.ResponseWriter, r *http.Request) {
	var (
		repo         repository.ReplyRepo = crud.NewReplyRepoCRUD()
		rid          uint64
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
	if rid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if updatedReply, status, err = repo.FindByID(rid); err != nil {
		response.Error(w, status, err)
		return
	}
	if input.Content != "" {
		updatedReply.Content = input.Content
	}
	if err = repo.Update(updatedReply); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been updated", nil)
}

func DeleteReply(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.ReplyRepo = crud.NewReplyRepoCRUD()
		rid    uint64
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

func DeleteAllRepliesForPost(pid uint64) error {
	var (
		repo repository.ReplyRepo = crud.NewReplyRepoCRUD()
		err  error
	)
	if err = repo.DeleteGroup(pid); err != nil {
		return err
	}
	return nil
}

func CountReplies(pid uint64) (replies string, err error) {
	var repo = crud.NewReplyRepoCRUD()
	if replies, err = repo.CountReplies(pid); err != nil {
		return "0", err
	}
	return replies, nil
}
