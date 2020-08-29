package middleware

import (
	"database/sql"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/database"
	"net/http"
)

type Middlewares func(http.HandlerFunc) http.HandlerFunc

//Chain function takes in multiple middleware functions,
//and combines them, to avoid spaghetti code.
func Chain(h http.HandlerFunc, m ...Middlewares) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}
	wrapped := h
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped
}

func checkUserRole(id uint64) (int, int, error) {
	var (
		role   int
		db     *sql.DB
		um     repository.UserRepo
		status int
		err    error
	)
	status = http.StatusInternalServerError
	if db, err = database.Connect(); err != nil {
		return 0, status, err
	}
	defer db.Close()
	if um, err = crud.NewUserModel(db); err != nil {
		return 0, status, err
	}
	if role, status, err = um.GetRole(id); err != nil {
		return role, status, err
	}
	return role, status, nil
}
