package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	models "forum/api/models/post"
	"forum/api/response"
	"forum/config"
	"forum/database"
	"net/http"
	"strconv"
	"time"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	db, dbErr := database.Connect()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	defer db.Close()
	pm, pmErr := models.NewPostModel(db)
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
	return
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, errors.New("invalid ID"))
		return
	}
	db, dbErr := database.Connect()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	defer db.Close()
	pm, pmErr := models.NewPostModel(db)
	if pmErr != nil {
		response.Error(w, http.StatusInternalServerError, pmErr)
		return
	}
	user, status, err := pm.FindByID(ID)
	if err != nil {
		response.Error(w, status, err)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, user)
	return
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
	db, dbErr := database.Connect()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	defer db.Close()
	pm, pmErr := models.NewPostModel(db)
	if pmErr != nil {
		response.Error(w, http.StatusInternalServerError, pmErr)
		return
	}
	post := models.Post{
		Name:     input.Description,
		Content:  input.Content,
		By:       uid,
		Category: 0,
		Created:  time.Now(),
		Updated:  time.Now(),
		Rating:   0,
	}

	createErr := pm.Create(&post)
	if createErr != nil {
		response.Error(w, http.StatusInternalServerError, createErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, "post has been created", nil)
	return
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("description")
	content := r.FormValue("content")
	ID, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	db, dbErr := database.Connect()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	defer db.Close()
	pm, pmErr := models.NewPostModel(db)
	if pmErr != nil {
		response.Error(w, http.StatusInternalServerError, pmErr)
		return
	}
	updatePost, status, findErr := pm.FindByID(ID)
	if findErr != nil {
		response.Error(w, status, findErr)
		return
	}
	updatePost.Updated = time.Now()
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
	return
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("invalid ID parameter"))
		return
	}
	db, dbErr := database.Connect()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	defer db.Close()
	pm, pmErr := models.NewPostModel(db)
	if pmErr != nil {
		response.Error(w, http.StatusInternalServerError, pmErr)
		return
	}
	if status, deleteErr := pm.Delete(ID); deleteErr != nil {
		response.Error(w, status, deleteErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("post has been deleted"), nil)
	return
}
