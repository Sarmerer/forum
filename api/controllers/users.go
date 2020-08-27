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
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	defer db.Close()
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
	return
}

//GetUser gets a specified user from the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64)
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
	return
}

//UpdateUser updates info about the user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	db, dbErr := database.Connect()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	defer db.Close()
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
	updatedUser, status, uErr := um.FindByID(ID)
	if uErr != nil {
		response.Error(w, status, uErr)
		return
	}
	if name != "" {
		updatedUser.Name = name
	}
	if status, updateErr := um.Update(updatedUser); updateErr != nil {
		response.Error(w, status, updateErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, "user has been updated", nil)
	return
}

//DeleteUser deletes a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64)
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
	return
}
