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
)

// Credentials struct models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"user_password" db:"user_password"`
	Username string `json:"user_name" db:"user_name"`
	Email    string `json:"user_email" db:"user_email"`
}

//SignIn signs the user in if exists
func SignIn(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	if login == "admin" && password == "root" {
		cookie, err := r.Cookie("sessionID")
		if err != http.ErrNoCookie {
			response.InternalError(w)
		}
		cookie = generateCookie()
		http.SetCookie(w, cookie)
	} else {
		response.Error(w, http.StatusBadRequest, errors.New("wrong login or password"))
	}
}

//SignUp authorizes new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	db, dbErr := database.Connect()
	um, umErr := models.NewUserModel(db)
	if dbErr != nil || umErr != nil {
		response.InternalError(w)
		return
	}
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{}
	decodeErr := json.NewDecoder(r.Body).Decode(creds)
	if decodeErr != nil {
		// If there is something wrong with the request body, return a 400 status
		response.BadRequest(w)
		return
	}
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, hashErr := security.Hash(creds.Password)
	if hashErr != nil {
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
	created, errText := um.Create(&user)
	//TO-DO: improve error check
	if !created {
		// If there is any issue with inserting into the database, return a 500 error
		response.Error(w, http.StatusInternalServerError, errors.New(errText))
		return
	}
	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err == http.ErrNoCookie {
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
