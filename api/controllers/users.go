package controllers

import (
	"errors"
	"strconv"

	"forum/api/controllers/helpers"
	models "forum/api/models/user"
	"forum/api/response"
	"forum/config"

	"net/http"
)

//GetUsers gets all users from the database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var (
		um    *models.UserModel
		users []models.User
		err   error
	)
	if um, err = helpers.NewUserModel(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if users, err = um.FindAll(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, users)
}

//GetUser gets a specified user from the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	var (
		uid    uint64
		um     *models.UserModel
		user   *models.User
		status int
		err    error
	)
	if uid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if um, err = helpers.NewUserModel(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if user, status, err = um.FindByID(uid); err != nil {
		response.Error(w, status, err)
		return
	}

	response.JSON(w, config.StatusSuccess, http.StatusOK, nil, user)
}

//UpdateUser updates info about the user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		name        string
		uid         uint64
		um          *models.UserModel
		updatedUser *models.User
		status      int
		err         error
	)
	name = r.FormValue("name")
	if uid, err = strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if um, err = helpers.NewUserModel(); err != nil {
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

	response.JSON(w, config.StatusSuccess, http.StatusOK, "user has been updated", nil)
}

//DeleteUser deletes a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var (
		uid    uint64
		um     *models.UserModel
		status int
		err    error
	)
	if uid, err = strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64); err != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("invalid ID parameter"))
		return
	}

	if um, err = helpers.NewUserModel(); err != nil {
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

	response.JSON(w, config.StatusSuccess, http.StatusOK, "user has been deleted", nil)
}
