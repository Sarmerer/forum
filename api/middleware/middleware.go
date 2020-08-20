package middleware

import (
	"errors"
	"forum/api/response"
	"forum/api/session"
	"forum/api/utils"
	"forum/config"
	"log"
	"net/http"
)

func SetupHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		//CORS handler
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next(w, r)
	}
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var method string
		switch r.Method {
		case http.MethodGet:
			method = utils.GET(r.Method)
		case http.MethodPost:
			method = utils.POST(r.Method)
		case http.MethodPut:
			method = utils.PUT(r.Method)
		case http.MethodDelete:
			method = utils.DELETE(r.Method)
		case http.MethodOptions:
			method = utils.OPTIONS(r.Method)
		default:
			method = utils.Default(r.Method)
		}
		log.Printf("\t%s |%s %s", r.Host, method, r.URL.Path)
		next(w, r)
	}
}

func AllowedMethods(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			response.Error(w, http.StatusMethodNotAllowed, errors.New("wrong method"))
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
