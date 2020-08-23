package middleware

import (
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
		if _, exists := cache.Sessions.Get(cookie.Value); !exists {
			response.Error(w, http.StatusUnauthorized, errors.New("user not authorized"))
			return
		}
		next(w, r)
	}
}

func SelfActionOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ID, err := strconv.ParseInt(r.URL.Query().Get("ID"), 10, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		cookie, cookieErr := r.Cookie(config.SessionCookieName)
		if cookieErr != nil {
			response.Error(w, http.StatusInternalServerError, cookieErr)
		}
		if item, sessionExists := cache.Sessions.Get(cookie.Value); sessionExists {
			moderator, modErr := checkUserRole(item.Belongs)
			if modErr != nil {
				response.Error(w, http.StatusInternalServerError, modErr)
				return
			}
			if item.Belongs != ID && !moderator {
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
