package controllers

import (
	"encoding/json"
	"errors"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/api/utils"

	"net/http"
)

//GetUsers gets all users from the database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var (
		repo  repository.UserRepo = crud.NewUserRepoCRUD()
		users []models.User
		err   error
	)
	if users, err = repo.FindAll(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, nil, users)
}

//GetUser gets a specified user from the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		uid    int64
		user   *models.User
		status int
		err    error
	)
	if uid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if user, status, err = repo.FindByID(uid); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, nil, user)
}

//UpdateUser updates info about the user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		repo        repository.UserRepo = crud.NewUserRepoCRUD()
		name        string
		uid         int64
		updatedUser *models.User
		status      int
		err         error
	)
	name = r.FormValue("name")
	if uid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if !requestorIsAccountOwner(r, uid) {
		response.Error(w, http.StatusForbidden, errors.New("this account doesn't belong to you"))
		return
	}

	if updatedUser, status, err = repo.FindByID(uid); err != nil {
		response.Error(w, status, err)
		return
	}
	if name != "" {
		updatedUser.DisplayName = name
	}
	if status, err = repo.Update(updatedUser); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, "user has been updated", nil)
}

//DeleteUser deletes a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		input  models.InputID
		status int
		err    error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if _, status, err = repo.FindByID(input.ID); err != nil {
		response.Error(w, status, err)
		return
	}
	if status, err = repo.Delete(input.ID); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, "user has been deleted", nil)
}
