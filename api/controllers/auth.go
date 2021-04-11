package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/sarmerer/forum/api/OAuth"
	"github.com/sarmerer/forum/api/services/emailverification"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"
)

func OAuthHandler(w http.ResponseWriter, r *http.Request) {
	var (
		providerName string         = strings.ToLower(r.FormValue("provider"))
		provider     OAuth.Provider = OAuth.Providers[providerName]
		cookie       string
		sessionID    string
		user         *models.User
		status       int
		err          error
	)

	if _, ok := OAuth.Providers[providerName]; !ok {
		response.Error(w, http.StatusBadRequest, errors.New("unknown provider"))
		return
	}
	cookie, sessionID = generateCookie(r.Cookie(config.SessionCookieName))

	if user, status, err = provider.Auth(r.URL.Query(), sessionID); err != nil {
		response.Error(w, status, err)
		return
	}

	w.Header().Set("Set-Cookie", cookie)
	response.Success(w, "user has been created", user)
}

// SignIn verifies user credinentials with database.
//
// Returns User model on success
func SignIn(w http.ResponseWriter, r *http.Request) {
	var (
		repo         repository.UserRepo = crud.NewUserRepoCRUD()
		input        models.InputUserSignIn
		user         *models.User
		userPassword string
		cookie       string
		newSessionID string
		status       int
		err          error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if user, status, err = repo.FindByLoginOrEmail([]string{input.Login}); err != nil {
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
	cookie, newSessionID = generateCookie(r.Cookie(config.SessionCookieName))
	if err = repo.UpdateSession(user.ID, newSessionID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Set-Cookie", cookie)
	response.Success(w, "user is logged in", user)
}

// SignUp creates new user record in database
//
// Returns User model on success
func SignUp(w http.ResponseWriter, r *http.Request) {
	var (
		repo           repository.UserRepo = crud.NewUserRepoCRUD()
		input          models.InputUserSignUp
		hashedPassword string
		admintToken    string = os.Getenv("ADMIN_TOKEN")
		verifiedState  bool
		newUser        *models.User
		status         int
		err            error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = input.Validate(); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if hashedPassword, err = hash(input.Password); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if !config.RequireEmailVerification {
		verifiedState = true
	}

	user := models.User{
		Username:  input.Login,
		Password:  hashedPassword,
		Email:     input.Email,
		Avatar:    url.PathEscape(path.Join("https://avatars.dicebear.com/api/male/", input.Login, ".svg")),
		Alias:     input.Login,
		SessionID: "",
		Role:      config.RoleUser,
		Verified:  verifiedState,
	}

	// This line compares environment variable with name ADMMIN_TOKEN,
	// if it exisits, and adminTolken, passed from user. If they match -
	// user gets an admin role, if not, user gets standard role
	if admintToken != "" && input.AdminToken == admintToken {
		user.Role = config.RoleAdmin
	}

	if newUser, status, err = repo.Create(&user); err != nil {
		response.Error(w, status, err)
		return
	}

	if config.RequireEmailVerification {
		if err = emailverification.Manager.SendVerificationEmail(newUser); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		response.Success(w, "email sent", nil)
		return
	}
	response.Success(w, "user created", nil)
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

func SendVerification(w http.ResponseWriter, r *http.Request) {
	var (
		email  string = r.FormValue("email")
		status int
		err    error
	)
	if email == "" {
		response.Error(w, http.StatusBadRequest, errors.New("no email provided"))
		return
	}
	if status, err = emailverification.Manager.ResendVerificationEmail(email); err != nil {
		response.Error(w, status, err)
		return
	}
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var (
		repo         repository.UserRepo = crud.NewUserRepoCRUD()
		code         string              = r.URL.Query().Get("code")
		email        string              = r.URL.Query().Get("email")
		user         *models.User
		cookie       string
		newSessionID string
		status       int
		err          error
	)
	if code == "" || email == "" {
		response.Error(w, http.StatusBadRequest, errors.New("no code or email present"))
		return
	}
	if status, err = emailverification.Manager.Verify(email, code); err != nil {
		response.Error(w, status, err)
		return
	}

	if user, status, err = repo.FindUnverifiedByEmail(email); err != nil {
		response.Error(w, status, err)
		return
	}

	if err = repo.Verify(user.ID, true); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	cookie, newSessionID = generateCookie(r.Cookie(config.SessionCookieName))

	if err = repo.UpdateSession(user.ID, newSessionID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Set-Cookie", cookie)
	response.Success(w, "user has been verified", user)
}

// Me function returns a user with an id from request context
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
