package controllers

import (
	"errors"
	"fmt"
	"log"

	"forum/api/entities"
	"forum/api/models"
	"forum/api/response"
	"forum/api/utils"
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
	db, _ := database.Connect()
	defer db.Close()
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
	ID, parseErr := utils.ParseURLInt(r.URL.Path, "/users/update/")
	if parseErr != nil {
		log.Println("Failed to parse, ", parseErr)
		response.BadRequest(w)
		return
	}
	updatedUser := &entities.User{
		ID:   int(ID),
		Name: r.FormValue("name"),
	}
	// err := json.NewDecoder(r.Body).Decode(user)
	// if err != nil {
	// 	// If there is something wrong with the request body, return a 400 status
	// 	response.Error(w, http.StatusBadRequest, errors.New("bad request"))
	// 	return
	// }
	updated := um.Update(updatedUser)
	if !updated {
		log.Println("Failed to update user ", updatedUser.Name)
		response.Error(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("update user ", updatedUser), nil)
}

//DeleteUser deletes a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURLInt(r.URL.Path, "/users/delete/")
	if err != nil {
		response.BadRequest(w)
		return
	}
	db, dbErr := database.Connect()
	defer db.Close()
	um, umErr := models.NewUserModel(db)
	if dbErr != nil || umErr != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	deleteErr := um.Delete(ID)
	if deleteErr != nil {
		// If there is any issue with deleting a user, return a 500 error
		response.Error(w, http.StatusInternalServerError, deleteErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("delete user ", ID), nil)
}
