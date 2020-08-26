package controllers

import (
	"errors"
	"fmt"
	"strconv"

	models "forum/api/models/user"
	"forum/api/response"
	"forum/config"
	"forum/database"

	"net/http"
)

//GetUsers gets all users from the database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, dbErr := database.Connect()
	defer db.Close()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	um, umErr := models.NewUserModel(db)
	if umErr != nil {
		response.Error(w, http.StatusInternalServerError, umErr)
		return
	}
	users, err := um.FindAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, users)
}

//GetUser gets a specified user from the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	db, dbErr := database.Connect()
	defer db.Close()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	um, umErr := models.NewUserModel(db)
	if umErr != nil {
		response.Error(w, http.StatusInternalServerError, umErr)
		return
	}
	user, status, err := um.FindByID(uid)
	if err != nil {
		response.Error(w, status, err)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, user)
}

//UpdateUser updates info about the user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db, dbErr := database.Connect()
	defer db.Close()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	um, umErr := models.NewUserModel(db)
	if umErr != nil {
		response.Error(w, http.StatusInternalServerError, umErr)
		return
	}
	ID, err := strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	updatedUser := &models.User{
		ID:   ID,
		Name: r.FormValue("name"),
	}
	if status, updateErr := um.Update(updatedUser); updateErr != nil {
		response.Error(w, status, updateErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, "user has been updated", nil)
}

//DeleteUser deletes a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("invalid ID parameter"))
		return
	}
	db, dbErr := database.Connect()
	defer db.Close()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	um, umErr := models.NewUserModel(db)
	if umErr != nil {
		response.Error(w, http.StatusInternalServerError, umErr)
		return
	}
	if status, deleteErr := um.Delete(ID); deleteErr != nil {
		response.Error(w, status, deleteErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("delete user ", ID), nil)
}
