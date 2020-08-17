package controllers

import (
	"encoding/json"
	"errors"
	"fmt"

	"forum/api/entities"
	"forum/api/models"
	"forum/api/response"
	"forum/api/utils"
	"forum/database"

	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, _ := database.Connect()
	um, _ := models.NewUserModel(db)
	users, err := um.FindAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("inernal server error"))
	}
	fmt.Println(users)
	w.Write([]byte(fmt.Sprint("users", users)))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	db, _ := database.Connect()
	um, _ := models.NewUserModel(db)
	var user entities.User
	ID, err := utils.ParseURL(r.URL.Path, "/users/")
	if err != nil {
		response.Error(w, http.StatusBadRequest, errors.New("bad request"))
		return
	}
	user, err = um.Find(int(ID.(int64)))
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
	w.Write([]byte(fmt.Sprint("get user ", ID)))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db, dbErr := database.Connect()
	um, umErr := models.NewUserModel(db)
	if dbErr != nil || umErr != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	// We're probably gonna need this ID to autofill update fields for admin (?)
	// ID, err := utils.ParseURL(r.URL.Path, "/users/update/")
	// if err != nil {
	// 	response.Error(w, http.StatusBadRequest, errors.New("bad request"))
	// 	return
	// }
	// Parse and decode the request body into a new `Credentials` instance
	user := &entities.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		response.Error(w, http.StatusBadRequest, errors.New("bad request"))
		return
	}
	updated := um.Update(user)
	if !updated {
		// If there is any issue with inserting into the database, return a 500 error
		response.Error(w, http.StatusInternalServerError, errors.New("Internal server error"))
		return
	}
	w.Write([]byte(fmt.Sprint("update user ", user)))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ID, err := utils.ParseURL(r.URL.Path, "/users/delete/")
	if err != nil {
		response.Error(w, http.StatusBadRequest, errors.New("bad request"))
		return
	}
	w.Write([]byte(fmt.Sprint("delete user ", ID)))
}
