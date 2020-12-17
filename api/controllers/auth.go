package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"
)

// LogIn verifies user credinentials with database.
//
// Returns User model on success
func LogIn(w http.ResponseWriter, r *http.Request) {
	var (
		repo         repository.UserRepo = crud.NewUserRepoCRUD()
		input        models.InputUserSignIn
		user         *models.User
		userPassword string
		cookie       string
		newUUID      string
		status       int
		err          error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if user, status, err = repo.FindByLoginOrEmail(input.Login); err != nil {
		response.Error(w, status, err)
		return
	}

	if userPassword, status, err = repo.GetPassword(user.ID); err != nil {
		response.Error(w, status, err)
		return
	}

	if err = verifyPassword(userPassword, input.Password); err != nil {
		response.Error(w, http.StatusBadRequest, errors.New("wrong login or password"))
		return
	}
	cookie, newUUID = generateCookie(r.Cookie(config.SessionCookieName))
	if err = repo.UpdateSession(user.ID, newUUID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Set-Cookie", cookie)
	response.Success(w, fmt.Sprint("user is logged in"), user)
}

// SignUp creates new user record in database
//
// Returns User model on success
func SignUp(w http.ResponseWriter, r *http.Request) {
	var (
		repo           repository.UserRepo = crud.NewUserRepoCRUD()
		input          models.InputUserSignUp
		hashedPassword string
		cookie         string
		role           int    = 0
		admintToken    string = os.Getenv("ADMIN_TOKEN")
		newSessionID   string
		newUser        *models.User
		status         int
		err            error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	// This line compares environment variable with name ADMMIN_TOKEN,
	// if it exisits, and adminTolken, passed from user. If they match -
	// user gets an admin role, if not, user gets standard role
	if admintToken != "" && input.Admin && input.AdminToken == admintToken {
		role = 2
	}

	if err = input.Validate(); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if hashedPassword, err = hash(input.Password); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	cookie, newSessionID = generateCookie(r.Cookie(config.SessionCookieName))
	user := models.User{
		Username:  input.Login,
		Password:  hashedPassword,
		Email:     input.Email,
		Avatar:    fmt.Sprintf("https://avatars.dicebear.com/api/male/%s.svg", input.Login),
		Alias:     input.Login,
		SessionID: newSessionID,
		Role:      role,
	}
	if newUser, status, err = repo.Create(&user); err != nil {
		response.Error(w, status, err)
		return
	}

	w.Header().Set("Set-Cookie", cookie)
	response.Success(w, "user has been created", newUser)
	return
}

// LogOut deletes user session id from database
// and destroys session cookie in his browser
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

// Me is an endpoint function, that helps to understand if user is authenticated.
// When client makes an api/me request, middleware checks if user is authenticated,
// if not, it returns 403 Forbidden http status
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
