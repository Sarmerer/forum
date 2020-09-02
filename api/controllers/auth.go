package controllers

import (
	"errors"
	"fmt"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/config"
	"net/http"
	"time"
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
	login := r.FormValue("login")
	password := r.FormValue("password")
	if user, status, err = repo.FindByNameOrEmail(login); err != nil {
		response.Error(w, status, err)
		return
	}
	if err = verifyPassword(user.Password, password); err != nil {
		response.Error(w, http.StatusBadRequest, errors.New("wrong login or password"))
		return
	}
	if cookie, err = r.Cookie(config.SessionCookieName); err == http.ErrNoCookie {
		cookie = generateCookie()
	} else {
		cookie.Expires = time.Now().Add(config.SessionExpiration)
	}
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
		Name:      login,
		Password:  string(hashedPassword),
		Email:     email,
		SessionID: "",
		Role:      config.RoleAdmin,
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
	cookie, _ = r.Cookie(config.SessionCookieName)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	response.Success(w, "user is logged out", nil)
	return
}

func Status(w http.ResponseWriter, r *http.Request) {
	// var (
	// 	cookie *http.Cookie
	// 	err    error
	// )
	// if cookie, err = r.Cookie(config.SessionCookieName); err == http.ErrNoCookie {

	// } else if err == nil {

	// }
}
