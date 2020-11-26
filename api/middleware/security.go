package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"
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
			if userCtx, status, err = repo.ValidateSession(cookie.Value); err != nil {
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

func AuthorizedOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if utils.GetUserFromCtx(r).Role < config.RoleUser {
			response.Error(w, http.StatusForbidden, errors.New("user not authorized"))
			return
		}
		next(w, r)
	}
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if utils.GetUserFromCtx(r).Role < config.RoleAdmin {
			response.Error(w, http.StatusForbidden, errors.New("permission denied"))
			return
		}
		next(w, r)
	}
}

func ModerOrAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if utils.GetUserFromCtx(r).Role < config.RoleModer {
			response.Error(w, http.StatusForbidden, errors.New("permission denied"))
			return
		}
		next(w, r)
	}
}
