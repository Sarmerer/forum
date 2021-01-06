package controllers

import (
	"encoding/json"
	"errors"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"

	"net/http"
)

//GetUsers returns all user records from database as an array
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

// FindUser allows to search user by various parameters.
// Currently supported:
//
// - id - returns a user with that id
func FindUser(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		input  models.InputFindUser
		user   *models.User
		status int
		err    error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	switch input.By {
	case "id":
		if user, status, err = repo.FindByID(input.ID); err != nil {
			response.Error(w, status, err)
			return
		}
	default:
		response.Error(w, http.StatusBadRequest, errors.New("unknown search type"))
		return
	}

	response.Success(w, nil, user)
}

//UpdateUser modifies user record in database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		repo        repository.UserRepo = crud.NewUserRepoCRUD()
		name        string
		userID      int64
		updatedUser *models.User
		status      int
		err         error
	)
	name = r.FormValue("name")
	if userID, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if !requestorIsEntityOwner(utils.GetUserFromCtx(r), userID) {
		response.Error(w, http.StatusForbidden, errors.New("this account doesn't belong to you"))
		return
	}

	if updatedUser, status, err = repo.FindByID(userID); err != nil {
		response.Error(w, status, err)
		return
	}
	if name != "" {
		updatedUser.Alias = name
	}
	if status, err = repo.Update(updatedUser); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, "user has been updated", nil)
}

//DeleteUser deletes a user from the database
// TODO delte all reactions of deleted user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		userID int64
		status int
		err    error
	)
	if userID, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if !requestorIsEntityOwner(utils.GetUserFromCtx(r), userID) {
		response.Error(w, http.StatusForbidden, errors.New("this account doesn't belong to you"))
		return
	}

	if _, status, err = repo.FindByID(userID); err != nil {
		response.Error(w, status, err)
		return
	}
	if status, err = repo.Delete(userID); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, "user has been deleted", nil)
}
