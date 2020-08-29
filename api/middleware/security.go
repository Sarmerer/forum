package middleware

import (
	"context"
	"database/sql"
	"errors"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/config"
	"forum/database"
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
		var (
			cookie *http.Cookie
			db     *sql.DB
			um     repository.UserRepo
			uid    uint64
			err    error
		)
		if cookie, err = r.Cookie(config.SessionCookieName); err == http.ErrNoCookie {
			response.Error(w, http.StatusUnauthorized, errors.New("user not authorized"))
			return
		}
		if db, err = database.Connect(); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		defer db.Close()
		if um, err = crud.NewUserModel(db); err != nil {
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
		queryUID, err := strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		requestorUID := r.Context().Value("uid").(uint64)
		if queryUID != requestorUID {
			role, status, modErr := checkUserRole(requestorUID)
			if modErr != nil {
				response.Error(w, status, modErr)
				return
			}
			if role == 0 {
				response.Error(w, http.StatusForbidden, errors.New("this account doesn't belong to you"))
				return
			}
		}
		next(w, r)
	}
}
