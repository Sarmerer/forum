package controllers

import (
	"fmt"

	"forum/api/entities"
	"forum/api/models"
	"forum/api/response"
	"forum/api/utils"
	"forum/config"
	"forum/database"

	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, _ := database.Connect()
	um, _ := models.NewUserModel(db)
	users, err := um.FindAll()
	if err != nil {
		response.InternalError(w)
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	db, _ := database.Connect()
	um, _ := models.NewUserModel(db)
	var user entities.User
	ID, err := utils.ParseURLInt(r.URL.Path, "/users/")
	if err != nil {
		response.BadRequest(w)
		return
	}
	user, err = um.Find(ID)
	if err != nil {
		panic(err)
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint("create user")))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURLInt(r.URL.Path, "/users/update/")
	if err != nil {
		response.BadRequest(w)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("update user ", ID), nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURLInt(r.URL.Path, "/users/delete/")
	if err != nil {
		response.BadRequest(w)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("delete user ", ID), nil)
}
