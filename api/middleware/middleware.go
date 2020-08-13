package middleware

import (
	"net/http"
)

func AllowedMethods(next http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			return
		}
		next(w, r)
	}
}
