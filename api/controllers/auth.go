package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/config"
	"net/http"
)

//SignIn signs the user in if exists
func SignIn(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		user   *models.User
		cookie *http.Cookie
		status int
		err    error
	)
	input := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}
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
	cookie = generateCookie(r.Cookie(config.SessionCookieName))
	if err = repo.UpdateSession(user.ID, cookie.Value); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, cookie)
	response.Success(w, fmt.Sprint("user is logged in"), nil)
}

//SignUp authorizes new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	var (
		repo           repository.UserRepo = crud.NewUserRepoCRUD()
		hashedPassword []byte
		status         int
		err            error
	)
	login := r.FormValue("login")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if login == "" || password == "" || email == "" {
		response.Error(w, http.StatusBadRequest, errors.New("empty login, email or password"))
		return
	}
	if hashedPassword, err = hash(password); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	user := models.User{
		Login:       login,
		Password:    string(hashedPassword),
		Email:       email,
		DisplayName: login,
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

func SignOut(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		uid    uint64              = r.Context().Value("uid").(uint64)
		cookie *http.Cookie
		err    error
	)
	if err = repo.UpdateSession(uid, ""); err != nil {
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
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		uid    uint64              = r.Context().Value("uid").(uint64)
		user   *models.User
		status int
		err    error
	)

	if user, status, err = repo.FindByID(uid); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, nil, user)
}
