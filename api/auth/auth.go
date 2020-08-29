package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/config"
	"forum/database"
	"net/http"
	"time"
)

//SignIn signs the user in if exists
func SignIn(w http.ResponseWriter, r *http.Request) {
	var (
		um     repository.UserRepo
		user   *models.User
		cookie *http.Cookie
		status int
		err    error
	)
	login := r.FormValue("login")
	password := r.FormValue("password")
	if login == "" || password == "" {
		response.Error(w, http.StatusBadRequest, errors.New("bad request"))
		return
	}
	if um, err = newUM(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}
	if user, status, err = um.FindByNameOrEmail(login); err != nil {
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
	if err = um.UpdateSession(user.ID, cookie.Value); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, cookie)
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("user is logged in"), nil)
}

//SignUp authorizes new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	var (
		hashedPassword []byte
		um             repository.UserRepo
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
	if um, err = newUM(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	user := models.User{
		Name:      login,
		Password:  string(hashedPassword),
		Email:     email,
		Nickname:  password,
		SessionID: "",
		Role:      0,
	}
	if status, err = um.Create(&user); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, "user has been created", nil)
	return
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	var (
		um  repository.UserRepo
		uid uint64
		err error
	)
	if um, err = newUM(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}
	uid = r.Context().Value("uid").(uint64)
	if err = um.UpdateSession(uid, ""); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	cookie, _ := r.Cookie(config.SessionCookieName)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	response.Success(w, "user is logged out", nil)
	return
}

func newUM() (um *crud.UserModel, err error) {
	var db *sql.DB
	if db, err = database.Connect(); err != nil {
		return
	}
	if um, err = crud.NewUserModel(db); err != nil {
		return
	}
	return
}
