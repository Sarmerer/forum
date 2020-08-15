package middleware

import (
	"errors"
	"forum/api/response"
	"forum/api/security"
	"log"
	"net/http"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s%s %s", r.Method, r.Host, r.RequestURI, r.Proto)
		next(w, r)
	}
}

func SetJSONType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func AllowedMethods(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			response.Error(w, http.StatusMethodNotAllowed, errors.New("Wrong Method"))
			return
		}
		next(w, r)
	}
}

func CheckUserAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sessionID")
		if err == http.ErrNoCookie {
			w.Write([]byte("user not authorized(from middleware.CheckUserAuth)"))
			return
		}
		sessionExists, err := security.ValidateSession(cookie.Value)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, errors.New("Internal server error"))
			return
		}
		if !sessionExists {
			response.Error(w, http.StatusUnauthorized, errors.New("user not authorized(from middleware.CheckUserAuth)"))
			return
		}
		next(w, r)
	}
}
