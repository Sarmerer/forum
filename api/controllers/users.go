package controllers

import (
	"fmt"

	"forum/api/entities"
	"forum/api/errors"
	"forum/api/models"
	"forum/api/utils"
	"forum/database"

	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, _ := database.Connect()
	um, _ := models.NewUserModel(db)
	users, err := um.FindAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(users)
	w.Write([]byte(fmt.Sprint(users)))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	db, _ := database.Connect()
	um, _ := models.NewUserModel(db)
	var user entities.User
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
	w.Write([]byte(fmt.Sprint("get user ", ID)))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint("create user")))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURL(r.URL.Path, "/users/update/")
	if err != nil {
		errors.HTTPErrorsHandler(http.StatusNotFound, w, r)
		return
	}
	w.Write([]byte(fmt.Sprint("update user ", ID)))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURL(r.URL.Path, "/users/delete/")
	if err != nil {
		errors.HTTPErrorsHandler(http.StatusNotFound, w, r)
		return
	}
	w.Write([]byte(fmt.Sprint("delete user ", ID)))
}
