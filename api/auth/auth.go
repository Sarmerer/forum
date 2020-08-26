package auth

import (
	"errors"
	"fmt"
	models "forum/api/models/user"
	"forum/api/response"
	"forum/config"
	"forum/database"
	"net/http"
	"time"
)

//SignIn signs the user in if exists
func SignIn(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	if login == "" || password == "" {
		response.Error(w, http.StatusBadRequest, errors.New("bad request"))
		return
	}
	db, dbErr := database.Connect()
	defer db.Close()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	um, umErr := models.NewUserModel(db)
	if umErr != nil {
		response.Error(w, http.StatusInternalServerError, umErr)
		return
	}
	user, uErr := um.FindByNameOrEmail(login)
	if uErr != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("inernal server error"))
		return
	}
	passErr := verifyPassword(user.Password, password)
	if passErr != nil {
		response.Error(w, http.StatusBadRequest, errors.New("wrong login or password"))
		return
	}
	cookie, cookieErr := r.Cookie(config.SessionCookieName)
	if cookieErr == http.ErrNoCookie {
		cookie = generateCookie()
	} else {
		cookie.Expires = time.Now().Add(config.SessionExpiration)
	}
	if err := um.UpdateSession(user.ID, cookie.Value); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, cookie)
	response.JSON(w, config.StatusSuccess, http.StatusOK, fmt.Sprint("user is logged in"), nil)
}

//SignUp authorizes new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if login == "" || password == "" || email == "" {
		response.Error(w, http.StatusBadRequest, errors.New("empty login, email or password"))
		return
	}
	hashedPassword, hashErr := hash(password)
	if hashErr != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("inernal server error"))
		return
	}
	db, dbErr := database.Connect()
	defer db.Close()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	um, umErr := models.NewUserModel(db)
	if umErr != nil {
		response.Error(w, http.StatusInternalServerError, umErr)
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
	// Next, insert the username, along with the hashed password into the database
	createErr := um.Create(&user)
	//TO-DO: improve error check
	if createErr != nil {
		// If there is any issue with inserting into the database, return a 500 error
		response.Error(w, http.StatusInternalServerError, createErr)
		return
	}
	response.Success(w, "user has been created", nil)
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	db, dbErr := database.Connect()
	defer db.Close()
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, dbErr)
		return
	}
	um, umErr := models.NewUserModel(db)
	if umErr != nil {
		response.Error(w, http.StatusInternalServerError, umErr)
		return
	}
	cookie, _ := r.Cookie(config.SessionCookieName)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	uid := r.Context().Value("uid").(uint64)
	if err := um.UpdateSession(uid, ""); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "user is logged out", nil)
}
