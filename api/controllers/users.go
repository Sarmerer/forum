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

func GetUser(w http.ResponseWriter, r *http.Request) {
	db, _ := database.Connect()
	um, _ := models.NewUserModel(db)
	var user entities.User
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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create user"))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update user"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete user"))
}
