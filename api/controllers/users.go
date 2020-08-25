package controllers

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"forum/api/cache"
	"forum/api/database"
	"forum/api/entities"
	"forum/api/models"
	"forum/api/response"
	"forum/config"

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
	ID, err := strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64)
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
	user, err := um.Find(ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
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
	updatedUser := &entities.User{
		ID:   ID,
		Name: r.FormValue("name"),
	}
	if updateErr := um.Update(updatedUser); updateErr != nil {
		log.Println("Failed to update user ", updatedUser.Name)
		response.Error(w, http.StatusInternalServerError, updateErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("update user ", updatedUser), nil)
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
	cookie, cookieErr := r.Cookie(config.SessionCookieName)
	if cookieErr != nil {
		response.Error(w, http.StatusInternalServerError, cookieErr)
		return
	}
	if cacheErr := cache.Sessions.Delete(cookie.Value); cacheErr != nil {
		response.Error(w, http.StatusInternalServerError, cacheErr)
		return
	}
	if deleteErr := um.Delete(ID); deleteErr != nil {
		response.Error(w, http.StatusInternalServerError, deleteErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("delete user ", ID), nil)
}
