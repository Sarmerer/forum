package controllers

import (
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/api/utils"
	"net/http"
)

func RatePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo repository.PostRepo = crud.NewPostRepoCRUD()
		uid  uint64              = r.Context().Value("uid").(uint64)
		mark int                 = -1
		pid  uint64
		err  error
	)
	if pid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if err = repo.RatePost(pid, uid, mark); err != nil {
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
