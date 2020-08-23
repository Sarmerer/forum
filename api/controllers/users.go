package controllers

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"forum/api/cache"
	"forum/api/entities"
	"forum/api/models"
	"forum/api/response"
	"forum/config"
	"forum/database"

	"net/http"
)

//GetUsers gets all users from the database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, _ := database.Connect()
	defer db.Close()
	um, _ := models.NewUserModel(db)
	users, err := um.FindAll()
	if err != nil {
		response.InternalError(w)
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, users)
}

//GetUser gets a specified user from the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()
	ID, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	um, err := models.NewUserModel(db)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	var user entities.User
	user, err = um.Find(ID)
	if err != nil {
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, user)
}

//UpdateUser updates info about the user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db, dbErr := database.Connect()
	defer db.Close()
	um, umErr := models.NewUserModel(db)
	if dbErr != nil || umErr != nil {
		log.Println("Failed to connect to the database")
		response.Error(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	ID, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	updatedUser := &entities.User{
		ID:   ID,
		Name: r.FormValue("name"),
	}
	// err := json.NewDecoder(r.Body).Decode(user)
	// if err != nil {
	// 	// If there is something wrong with the request body, return a 400 status
	// 	response.Error(w, http.StatusBadRequest, errors.New("bad request"))
	// 	return
	// }
	if updateErr := um.Update(updatedUser); updateErr != nil {
		log.Println("Failed to update user ", updatedUser.Name)
		response.Error(w, http.StatusInternalServerError, updateErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("update user ", updatedUser), nil)
}

//DeleteUser deletes a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("invalid ID parameter"))
		return
	}
	db, dbErr := database.Connect()
	defer db.Close()
	um, umErr := models.NewUserModel(db)
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
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
