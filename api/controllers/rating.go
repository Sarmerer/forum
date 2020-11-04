package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"
)

func RatePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.PostRepo = crud.NewPostRepoCRUD()
		userCtx models.UserCtx      = utils.GetUserFromCtx(r)
		input   models.InputRate
		result  models.Rating
		err     error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if err = repo.Rate(input.ID, userCtx.ID, input.Reaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if result.Rating, result.YourReaction, err = repo.GetRating(input.ID, userCtx.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "post has been rated", result)
}

func DeleteReactionsForPost(pid int64) error {
	var (
		repo repository.PostRepo = crud.NewPostRepoCRUD()
		err  error
	)
	if err = repo.DeleteAllReactions(pid); err != nil {
		return err
	}
	return nil
}

func RateComment(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.CommentRepo = crud.NewCommentRepoCRUD()
		userCtx models.UserCtx         = utils.GetUserFromCtx(r)
		input   models.InputRate
		result  models.Rating
		err     error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if err = repo.Rate(input.ID, userCtx.ID, input.Reaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if result.Rating, result.YourReaction, err = repo.GetRating(input.ID, userCtx.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "comment has been rated", result)
}