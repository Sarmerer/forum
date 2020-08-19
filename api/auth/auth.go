package auth

import (
	"encoding/json"
	"errors"
	"forum/api/entities"
	"forum/api/models"
	"forum/api/response"
	"forum/config"
	"forum/database"
	"log"
	"net/http"
)

// credentials struct models the structure of a user, both in the request body, and in the DB
type credentials struct {
	Username string `json:"user_name" db:"user_name"`
	Password string `json:"user_password" db:"user_password"`
	Email    string `json:"user_email" db:"user_email"`
}

//SignIn signs the user in if exists
func SignIn(w http.ResponseWriter, r *http.Request) {
	db, dbErr := database.Connect()
	um, umErr := models.NewUserModel(db)
	defer db.Close()
	if dbErr != nil || umErr != nil {
		response.InternalError(w)
		return
	}
	login := r.FormValue("login")
	password := r.FormValue("password")
	if login == "" || password == "" {
		response.BadRequest(w)
		return
	}
	creds := &credentials{Username: login, Password: password}
	user, uErr := um.FindByNameOrEmail(creds.Username)
	if uErr != nil {
		response.InternalError(w)
		return
	}
	passErr := verifyPassword(user.Password, creds.Password)
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
	user.SessionID = cookie.Value
	um.Update(&user)
	response.JSON(w, config.StatusSuccess, http.StatusOK, "user is logged in", nil)
}

//SignUp authorizes new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &credentials{}
	decodeErr := json.NewDecoder(r.Body).Decode(creds)
	if decodeErr != nil {
		// If there is something wrong with the request body, return a 400 status
		response.BadRequest(w)
		return
	}
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, hashErr := hash(creds.Password)
	if hashErr != nil {
		response.InternalError(w)
		return
	}
	db, dbErr := database.Connect()
	defer db.Close()
	um, umErr := models.NewUserModel(db)
	if dbErr != nil || umErr != nil {
		response.InternalError(w)
		return
	}
	user := entities.User{
		Name:      creds.Username,
		Password:  string(hashedPassword),
		Email:     creds.Email,
		Nickname:  creds.Username,
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
	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(config.SessionCookieName)
	if err == http.ErrNoCookie {
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
