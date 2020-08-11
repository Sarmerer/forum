package auth

import (
	"forum/api/errors"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/signin":
			signIn(w, r)
		case "/signup":
			signUp(w, r)
		case "/signout":
			signOut(w, r)
		default:
			errors.HTTPErrorsHandler(404, w, r)
		}
	default:
		errors.HTTPErrorsHandler(405, w, r)
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
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
func signUp(w http.ResponseWriter, r *http.Request) {
	//TODO: create new user
}

func signOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err == http.ErrNoCookie {
		return
	}
	delete(activeSessions, uuid.FromStringOrNil(cookie.Value))
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
