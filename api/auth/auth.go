package auth

import (
	"errors"
	"fmt"
	"forum/api/cache"
	"forum/api/entities"
	"forum/api/models"
	"forum/api/response"
	"forum/config"
	"log"
	"net/http"
)

//SignIn signs the user in if exists
func SignIn(w http.ResponseWriter, r *http.Request) {
	db, um, umErr := models.NewUserModel()
	defer db.Close()
	if umErr != nil {
		response.Error(w, http.StatusInternalServerError, umErr)
		return
	}
	login := r.FormValue("login")
	password := r.FormValue("password")
	if login == "" || password == "" {
		response.BadRequest(w)
		return
	}
	user, uErr := um.FindByNameOrEmail(login)
	if uErr != nil {
		response.InternalError(w)
		return
	}
	passErr := verifyPassword(user.Password, password)
	if passErr != nil {
		response.Error(w, http.StatusBadRequest, errors.New("wrong login or password"))
		return
	}
	cookie, cookieErr := r.Cookie(config.SessionCookieName)
	if cookieErr != http.ErrNoCookie && cookieErr != nil {
		log.Println("Failed to generate cookie: ", cookieErr)
		response.InternalError(w)
		return
	}
	cookie = generateCookie()
	http.SetCookie(w, cookie)
	cache.Sessions.Set(cookie.Value, &cache.Session{SessionID: cookie.Value, Belongs: user.ID}, 0)
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("user is logged in"), nil)
}

//SignUp authorizes new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if login == "" || password == "" || email == "" {
		response.BadRequest(w)
		return
	}
	hashedPassword, hashErr := hash(password)
	if hashErr != nil {
		response.InternalError(w)
		return
	}
	db, um, umErr := models.NewUserModel()
	defer db.Close()
	if umErr != nil {
		response.Error(w, http.StatusInternalServerError, umErr)
		return
	}
	user := entities.User{
		Name:      login,
		Password:  string(hashedPassword),
		Email:     email,
		Nickname:  password,
		SessionID: "",
		Role:      0,
	}
	// Next, insert the username, along with the hashed password into the database
	createErr := um.Create(&user)
	//TO-DO: improve error check
	if createErr != nil {
		// If there is any issue with inserting into the database, return a 500 error
		response.Error(w, http.StatusInternalServerError, createErr)
		return
	}
	response.JSON(w, config.StatusSuccess, http.StatusOK, "user has been created", nil)
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(config.SessionCookieName)
	if err == http.ErrNoCookie {
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	cache.Sessions.Delete(cookie.Value)
}
