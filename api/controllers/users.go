package controllers

import (
	"fmt"
	"forum/api/errors"
	"forum/utils"
	"net/http"
)

type User struct {
	ID int
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUser(w, r)
	case http.MethodPost:
		createUser(w, r)
	case http.MethodPut:
		updateUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		errors.HTTPErrorsHandler(http.StatusMethodNotAllowed, w, r)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURL(r.URL.Path, "/users/")
	if err != nil {
		errors.HTTPErrorsHandler(http.StatusNotFound, w, r)
		return
	}
	w.Write([]byte(fmt.Sprint("get user ", ID)))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create user"))
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update user"))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete user"))
}
