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
		err     error
	)
	input := struct {
		PID      int64 `json:"pid"`
		Reaction int   `json:"reaction"`
	}{}
	result := struct {
		Rating       int `json:"rating"`
		YourReaction int `json:"your_reaction"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if err = repo.RatePost(input.PID, userCtx.ID, input.Reaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if result.Rating, result.YourReaction, err = repo.GetRating(input.PID, userCtx.ID); err != nil {
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
