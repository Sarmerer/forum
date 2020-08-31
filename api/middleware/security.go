package middleware

import (
	"context"
	"errors"
	"forum/api/helpers"
	"forum/api/repository"
	"forum/api/response"
	"forum/config"
	"net/http"
	"os"
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
		var (
			cookie *http.Cookie
			um     repository.UserRepo
			uid    uint64
			err    error
		)
		if cookie, err = r.Cookie(config.SessionCookieName); err == http.ErrNoCookie {
			response.Error(w, http.StatusUnauthorized, errors.New("user not authorized"))
			return
		}
		if um, err = helpers.PrepareUserRepo(); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if uid, err = um.ValidateSession(cookie.Value); err != nil {
			response.Error(w, http.StatusUnauthorized, errors.New("user not authorized"))
			return
		}
		ctx := context.WithValue(r.Context(), "uid", uid)
		next(w, r.WithContext(ctx))
	}
}

func SelfActionOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			queryUID     uint64
			requestorUID uint64
			role         int
			um           repository.UserRepo
			status       int
			err          error
		)
		if queryUID, err = helpers.ParseID(r); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		requestorUID = r.Context().Value("uid").(uint64)
		if queryUID != requestorUID {
			if um, err = helpers.PrepareUserRepo(); err != nil {
				response.Error(w, http.StatusInternalServerError, err)
				return
			}
			if role, status, err = um.GetRole(requestorUID); err != nil {
				response.Error(w, status, err)
				return
			} else if role < config.RoleModerator {
				response.Error(w, http.StatusForbidden, errors.New("this account doesn't belong to you"))
				return
			}
		}
		next(w, r)
	}
}
