package controllers

import (
	"fmt"
	"forum/utils"
	"net/http"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPost(w, r)
	case http.MethodPost:
		createPost(w, r)
	case http.MethodPut:
		updatePost(w, r)
	case http.MethodDelete:
		deletePost(w, r)
	}
}

func getPost(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURL(r.URL.Path, "/posts/")
	if err != nil {
		utils.HTTPErrorsHandler(http.StatusNotFound, w, r)
		return
	}
	w.Write([]byte(fmt.Sprint("get post", ID)))
}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create post"))
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update post"))
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete post"))
}
