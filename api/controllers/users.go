package controllers

import (
	"errors"

	"forum/api/helpers"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/response"

	"net/http"
)

//GetUsers gets all users from the database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var (
		um    repository.UserRepo
		users []models.User
		err   error
	)
	if um, err = helpers.PrepareUserRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if users, err = um.FindAll(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, nil, users)
}

//GetUser gets a specified user from the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	var (
		uid    uint64
		um     repository.UserRepo
		user   *models.User
		status int
		err    error
	)
	if uid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if um, err = helpers.PrepareUserRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if user, status, err = um.FindByID(uid); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, nil, user)
}

//UpdateUser updates info about the user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		name        string
		uid         uint64
		um          repository.UserRepo
		updatedUser *models.User
		status      int
		err         error
	)
	name = r.FormValue("name")
	if uid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if um, err = helpers.PrepareUserRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if updatedUser, status, err = um.FindByID(uid); err != nil {
		response.Error(w, status, err)
		return
	}
	if name != "" {
		updatedUser.Name = name
	}
	if status, err = um.Update(updatedUser); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, "user has been updated", nil)
}

//DeleteUser deletes a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var (
		uid    uint64
		um     repository.UserRepo
		status int
		err    error
	)
	if uid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("invalid ID parameter"))
		return
	}

	if um, err = helpers.PrepareUserRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if _, status, err = um.FindByID(uid); err != nil {
		response.Error(w, status, err)
		return
	}
	if status, err = um.Delete(uid); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, "user has been deleted", nil)
}
