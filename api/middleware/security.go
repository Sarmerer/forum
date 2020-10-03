package middleware

import (
	"context"
	"errors"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/api/utils"
	"forum/config"
	"net/http"
)

func CheckUserAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			repo   repository.UserRepo = crud.NewUserRepoCRUD()
			cookie *http.Cookie
			uid    int64
			status int
			err    error
		)
		if cookie, err = r.Cookie(config.SessionCookieName); err != nil {
			if err != http.ErrNoCookie {
				response.Error(w, http.StatusBadRequest, err)
			}
			response.Error(w, http.StatusUnauthorized, errors.New("user not authorized"))
			return
		}
		if uid, status, err = repo.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		ctx := context.WithValue(r.Context(), "uid", uid)
		next(w, r.WithContext(ctx))
	}
}

func SelfActionOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			repo         repository.UserRepo = crud.NewUserRepoCRUD()
			requestorUID int64              = r.Context().Value("uid").(int64)
			queryUID     int64
			role         int
			status       int
			err          error
		)
		if queryUID, err = utils.ParseID(r); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		if queryUID != requestorUID {

			if role, status, err = repo.GetRole(requestorUID); err != nil {
				response.Error(w, status, err)
				return
			} else if role < config.RoleAdmin {
				response.Error(w, http.StatusForbidden, errors.New("this account doesn't belong to you"))
				return
			}
		}
		next(w, r)
	}
}
