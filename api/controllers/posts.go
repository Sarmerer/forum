package controllers

import (
	"encoding/json"
	"fmt"
	"forum/api/controllers/helpers"
	models "forum/api/models/post"
	"forum/api/response"
	"forum/config"
	"net/http"
	"time"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var (
		pm    *models.PostModel
		posts []models.Post
		err   error
	)
	if pm, err = helpers.NewPostModel(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if posts, err = pm.FindAll(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	var (
		pid    uint64
		pm     *models.PostModel
		post   *models.Post
		status int
		err    error
	)
	if pid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if pm, err = helpers.NewPostModel(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if post, status, err = pm.FindByID(pid); err != nil {
		response.Error(w, status, err)
		return
	}

	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, post)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var (
		author uint64
		pm     *models.PostModel
		post   models.Post
		err    error
	)
	author = r.Context().Value("uid").(uint64)
	input := struct {
		Description string `json:"description"`
		Content     string `json:"content"`
	}{}

	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if pm, err = helpers.NewPostModel(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	post = models.Post{
		Name:     input.Description,
		Content:  input.Content,
		By:       author,
		Category: 0,
		Created:  time.Now(),
		Updated:  time.Now(),
		Rating:   0,
	}
	if err = pm.Create(&post); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, config.StatusSuccess, http.StatusOK, "post has been created", nil)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var (
		name        string
		content     string
		pid         uint64
		pm          *models.PostModel
		updatedPost *models.Post
		status      int
		err         error
	)
	name = r.FormValue("description")
	content = r.FormValue("content")
	if pid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if pm, err = helpers.NewPostModel(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if status, err = helpers.PostExists(pid); err != nil {
		response.Error(w, status, err)
		return
	}

	updatedPost.Updated = time.Now()
	if name != "" {
		updatedPost.Name = name
	}
	if content != "" {
		updatedPost.Content = content
	}

	if err = pm.Update(updatedPost); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("post has been updated"), nil)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	var (
		pid    uint64
		pm     *models.PostModel
		status int
		err    error
	)
	if pid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if pm, err = helpers.NewPostModel(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if status, err = helpers.PostExists(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	if status, err = pm.Delete(pid); err != nil {
		response.Error(w, status, err)
		return
	}

	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("post has been deleted"), nil)
}
