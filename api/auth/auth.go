package auth

import (
	"encoding/json"
	"errors"
	"forum/api/entities"
	"forum/api/models"
	"forum/api/response"
	"forum/api/security"
	"forum/database"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Credentials struct models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"user_password" db:"user_password"`
	Username string `json:"user_name" db:"user_name"`
	Email    string `json:"user_email" db:"user_email"`
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if r.FormValue("password") == "admin" && r.FormValue("login") == "root" {
		if err == http.ErrNoCookie {
			err = nil
			cookie = &http.Cookie{Name: "sessionID", Value: ""}
		}
		cookie = generateCookie()
		http.SetCookie(w, cookie)
	} else {
		response.Error(w, http.StatusBadRequest, errors.New("Wrong login or password"))
	}
}

//SignUp authorizes new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	db, dbErr := database.Connect()
	um, umErr := models.NewUserModel(db)
	if dbErr != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("Internal server error"))
		return
	}
	if umErr != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("Internal server error"))
		return
	}
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		response.Error(w, http.StatusBadRequest, errors.New("Bad request"))
		return
	}
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := security.Hash(creds.Password)

	user := entities.User{
		Name:      creds.Username,
		Password:  string(hashedPassword),
		Email:     creds.Email,
		Nickname:  creds.Username,
		SessionID: "",
		Role:      0,
	}
	// Next, insert the username, along with the hashed password into the database
	if !um.Create(&user) {
		// If there is any issue with inserting into the database, return a 500 error
		response.Error(w, http.StatusInternalServerError, errors.New("Internal server error"))
		return
	}
	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err == http.ErrNoCookie {
		return
	}
	delete(activeSessions, uuid.FromStringOrNil(cookie.Value))
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
