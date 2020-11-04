package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"
)

//LogIn signs the user in if exists
func LogIn(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.UserRepo = crud.NewUserRepoCRUD()
		input   models.InputUserSignIn
		user    *models.User
		cookie  string
		newUUID string
		status  int
		err     error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if user, status, err = repo.FindByNameOrEmail(input.Login); err != nil {
		response.Error(w, status, err)
		return
	}
	if err = verifyPassword(user.Password, input.Password); err != nil {
		response.Error(w, http.StatusBadRequest, errors.New("wrong login or password"))
		return
	}
	cookie, newUUID = generateCookie(r.Cookie(config.SessionCookieName))
	if err = repo.UpdateSession(user.ID, newUUID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Set-Cookie", cookie)
	response.Success(w, fmt.Sprint("user is logged in"), nil)
}

//SignUp authorizes new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	var (
		repo           repository.UserRepo = crud.NewUserRepoCRUD()
		input          models.InputUserSignUp
		hashedPassword []byte
		status         int
		err            error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if hashedPassword, err = hash(input.Password); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	user := models.User{
		Login:       input.Login,
		Password:    string(hashedPassword),
		Email:       input.Email,
		Avatar:      "https://avatars.dicebear.com/api/male/" + input.Login + ".svg",
		DisplayName: input.Login,
		SessionID:   "",
		Role:        config.RoleAdmin,
	}
	if status, err = repo.Create(&user); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, "user has been created", nil)
	return
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.UserRepo = crud.NewUserRepoCRUD()
		userCtx models.UserCtx      = utils.GetUserFromCtx(r)
		cookie  *http.Cookie
		err     error
	)
	if err = repo.UpdateSession(userCtx.ID, ""); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	cookie = &http.Cookie{
		Name:     config.SessionCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	response.Success(w, "user is logged out", nil)
	return
}

func Me(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.UserRepo = crud.NewUserRepoCRUD()
		userCtx models.UserCtx      = utils.GetUserFromCtx(r)
		user    *models.User
		status  int
		err     error
	)
	if user, status, err = repo.FindByID(userCtx.ID); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, nil, user)
}
