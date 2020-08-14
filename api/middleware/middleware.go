package middleware

import (
	"net/http"
)

func AllowedMethods(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Write([]byte("405 Wrong Method(message from middleware.AllowedMethods)"))
			return
		}
		next(w, r)
	}
}
