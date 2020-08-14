package auth

import (
	"encoding/json"
	"fmt"
	"forum/api/entities"
	"forum/api/models"
	"forum/database"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
	bcrypt "golang.org/x/crypto/bcrypt"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println(formatRequest(r))
	cookie, err := r.Cookie("sessionID")
	if r.FormValue("password") == "admin" && r.FormValue("login") == "root" {
		if err == http.ErrNoCookie {
			err = nil
			cookie = &http.Cookie{Name: "sessionID", Value: ""}
		}
		cookie = generateCookie()
		http.SetCookie(w, cookie)
	} else {
		w.Write([]byte("Wrong login or password"))
	}
}

//SignUp authorizes new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	db, _ := database.Connect()
	um, _ := models.NewUserModel(db)
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
		return
	}
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

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
		w.WriteHeader(http.StatusInternalServerError)
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

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
