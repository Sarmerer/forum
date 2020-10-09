package controllers

import (
	"encoding/json"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"net/http"
)

func RatePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo repository.PostRepo = crud.NewPostRepoCRUD()
		uid  int64               = r.Context().Value("uid").(int64)
		err  error
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
	if err = repo.RatePost(input.PID, uid, input.Reaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if result.Rating, result.YourReaction, err = GetRating(input.PID, uid); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "post has been rated", result)
}

func GetRating(postID int64, uid int64) (int, int, error) {
	var (
		repo         repository.PostRepo = crud.NewPostRepoCRUD()
		rating       int
		yourReaction int
		err          error
	)
	if rating, yourReaction, err = repo.GetRating(postID, uid); err != nil {
		return 0, 0, err
	}
	return rating, yourReaction, nil
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
