package controllers

import (
	"fmt"

	"forum/api/entities"
	"forum/api/errors"
	"forum/api/models"
	"forum/api/utils"

	"net/http"
)

var um models.UserModel
var user entities.User

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
	user, err = um.Find(int(ID.(int64)))
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
	w.Write([]byte(fmt.Sprint(user)))
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
