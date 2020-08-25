package middleware

import (
	"context"
	"errors"
	"forum/api/cache"
	"forum/api/response"
	"forum/config"
	"net/http"
	"os"
	"strconv"
)

func CheckAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("API_KEY")
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
		session, exists := cache.Sessions.Get(cookie.Value)
		if !exists {
			response.Error(w, http.StatusUnauthorized, errors.New("user not authorized"))
			return
		}
		ctx := context.WithValue(r.Context(), "uid", session.Belongs)
		next(w, r.WithContext(ctx))
	}
}

func SelfActionOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, err := strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		cookie, cookieErr := r.Cookie(config.SessionCookieName)
		if cookieErr != nil {
			response.Error(w, http.StatusInternalServerError, cookieErr)
		}
		if session, exists := cache.Sessions.Get(cookie.Value); exists {
			moderator, modErr := checkUserRole(session.Belongs)
			if modErr != nil {
				response.Error(w, http.StatusInternalServerError, modErr)
				return
			}
			if session.Belongs != uid && !moderator {
				response.Error(w, http.StatusForbidden, errors.New("you can not delete someone else's account"))
				return
			}
		} else {
			response.Error(w, http.StatusForbidden, errors.New("user not authorized"))
			return
		}
		next(w, r)
	}
}
