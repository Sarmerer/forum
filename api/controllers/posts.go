package controllers

import (
	"fmt"
	"forum/api/errors"
	"forum/api/utils"
	"net/http"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint("get all posts")))
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURL(r.URL.Path, "/posts/")
	if err != nil {
		errors.HTTPErrorsHandler(http.StatusNotFound, w, r)
		return
	}
	w.Write([]byte(fmt.Sprint("get post", ID)))
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint("create post")))
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURL(r.URL.Path, "/posts/update/")
	if err != nil {
		errors.HTTPErrorsHandler(http.StatusNotFound, w, r)
		return
	}
	w.Write([]byte(fmt.Sprint("update post ", ID)))
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURL(r.URL.Path, "/posts/delete/")
	if err != nil {
		errors.HTTPErrorsHandler(http.StatusNotFound, w, r)
		return
	}
	w.Write([]byte(fmt.Sprint("delete post ", ID)))
}
