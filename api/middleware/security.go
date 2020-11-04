package middleware

import (
	"context"
	"errors"
	"forum/api/config"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"net/http"
)

func SetContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			repo    repository.UserRepo = crud.NewUserRepoCRUD()
			cookie  *http.Cookie
			userCtx models.UserCtx
			status  int
			ctx     context.Context
			err     error
		)
		if cookie, err = r.Cookie(config.SessionCookieName); err != nil {
			if err != http.ErrNoCookie {
				response.Error(w, http.StatusBadRequest, err)
			}
			userCtx = models.UserCtx{ID: -1, Role: -1}
			ctx = context.WithValue(r.Context(), config.UserCtxVarName, userCtx)
			next(w, r.WithContext(ctx))
		} else {
			if userCtx.ID, userCtx.Role, status, err = repo.ValidateSession(cookie.Value); err != nil {
				if err != nil && status != http.StatusUnauthorized {
					response.Error(w, status, err)
					return
				}
			}
			ctx = context.WithValue(r.Context(), config.UserCtxVarName, userCtx)
			next(w, r.WithContext(ctx))
		}
	}
}

func CheckUserAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			repo   repository.UserRepo = crud.NewUserRepoCRUD()
			cookie *http.Cookie
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
		if _, _, status, err = repo.ValidateSession(cookie.Value); err != nil {
			response.Error(w, status, err)
			return
		}
		next(w, r)
	}
}

// func RequestorIsAccountOwner(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var (
// 			repo         repository.UserRepo = crud.NewUserRepoCRUD()
// 			requestorUID int64               = r.Context().Value(config.UserCtxVarName).(int64)
// 			queryUID     int64
// 			role         int
// 			status       int
// 			err          error
// 		)
// 		if queryUID, err = utils.ParseID(r); err != nil {
// 			response.Error(w, http.StatusBadRequest, err)
// 			return
// 		}
// 		if queryUID != requestorUID {

// 			if role, status, err = repo.GetRole(requestorUID); err != nil {
// 				response.Error(w, status, err)
// 				return
// 			} else if role < config.RoleAdmin {
// 				response.Error(w, http.StatusForbidden, errors.New("this account doesn't belong to you"))
// 				return
// 			}
// 		}
// 		next(w, r)
// 	}
// }
