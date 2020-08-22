package middleware

import (
	"errors"
	"fmt"
	"forum/api/response"
	"forum/api/session"
	"forum/config"
	"net/http"
	"os"
)

func CheckAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("API_KEY")
		fmt.Println(key)
		if key == "" || key != os.Getenv("API_KEY") {
			response.Error(w, http.StatusForbidden, errors.New("cannot access API without a valid API key"))
			return
		}
		next(w, r)
	}
}

func CheckUserAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(config.SessionCookieName)
		if err == http.ErrNoCookie {
			response.Error(w, http.StatusUnauthorized, errors.New("user not authorized"))
			return
		}
		sessionExists, err := session.Validate(cookie.Value)
		if err != nil {
			response.InternalError(w)
			return
		}
		if !sessionExists {
			response.Error(w, http.StatusUnauthorized, errors.New("user not authorized"))
			return
		}
		next(w, r)
	}
}
