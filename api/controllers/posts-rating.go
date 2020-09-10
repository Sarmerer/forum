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
		uid  uint64              = r.Context().Value("uid").(uint64)
		err  error
	)
	input := struct {
		Pid      uint64 `json:"pid"`
		Reaction int    `json:"reaction"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if err = repo.RatePost(input.Pid, uid, input.Reaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "post has been rated", nil)
}

func GetRating(postID uint64) (int, error) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		rating int
		err    error
	)
	if rating, err = repo.GetRating(postID); err != nil {
		return 0, err
	}
	return rating, nil
}
