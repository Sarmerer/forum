package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/api/entities"
	"forum/api/models"
	"forum/api/response"
	"forum/config"
	"net/http"
	"strconv"
	"time"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	pm, pmErr := models.NewPostModel()
	defer pm.DB.Close()
	if pmErr != nil {
		response.Error(w, http.StatusInternalServerError, pmErr)
		return
	}
	posts, err := pm.FindAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("invalid ID parameter"))
		return
	}
	pm, pmErr := models.NewPostModel()
	defer pm.DB.Close()
	if pmErr != nil {
		response.Error(w, http.StatusInternalServerError, pmErr)
		return
	}
	user, err := pm.Find(ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, user)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	input := struct {
		Description string `json:"description"`
		Content     string `json:"content"`
	}{}
	jErr := json.NewDecoder(r.Body).Decode(&input)
	if jErr != nil {
		response.Error(w, http.StatusBadRequest, jErr)
		return
	}
	uid := r.Context().Value("uid").(uint64)
	pm, pmErr := models.NewPostModel()
	defer pm.DB.Close()
	if pmErr != nil {
		response.Error(w, http.StatusInternalServerError, pmErr)
		return
	}
	post := entities.Post{
		Name:     input.Description,
		Content:  input.Content,
		By:       uid,
		Category: 0,
		Created:  time.Now(),
		Updated:  time.Now(),
		Rating:   0,
	}
	// Next, insert the username, along with the hashed password into the database
	createErr := pm.Create(&post)
	if createErr != nil {
		response.Error(w, http.StatusInternalServerError, createErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, "post has been created", nil)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("description")
	content := r.FormValue("content")
	ID, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	pm, pmErr := models.NewPostModel()
	defer pm.DB.Close()
	if pmErr != nil {
		response.Error(w, http.StatusInternalServerError, pmErr)
		return
	}
	updatePost, findErr := pm.Find(ID)
	updatePost.Updated = time.Now()
	if findErr != nil {
		response.Error(w, http.StatusInternalServerError, findErr)
	}
	if name != "" {
		updatePost.Name = name
	}
	if content != "" {
		updatePost.Content = content
	}
	if updateErr := pm.Update(updatePost); updateErr != nil {
		response.Error(w, http.StatusInternalServerError, updateErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("post has been updated"), nil)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("invalid ID parameter"))
		return
	}
	pm, pmErr := models.NewPostModel()
	defer pm.DB.Close()
	if pmErr != nil {
		response.Error(w, http.StatusInternalServerError, pmErr)
		return
	}
	if deleteErr := pm.Delete(ID); deleteErr != nil {
		response.Error(w, http.StatusInternalServerError, deleteErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("post has been deleted"), nil)
}
