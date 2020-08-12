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

func GetUsers(w http.ResponseWriter, r *http.Request) {
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURL(r.URL.Path, "/user/")
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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create user"))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update user"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete user"))
}
