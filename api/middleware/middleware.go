package middleware

import (
	"errors"
	"forum/api/response"
	"forum/api/utils"
	"log"
	"net/http"
)

func SetHeaders(next http.HandlerFunc) http.HandlerFunc {
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
		log.Printf("\t%s |%s %s", r.Host, utils.Paint(r.Method), r.URL.Path)
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
