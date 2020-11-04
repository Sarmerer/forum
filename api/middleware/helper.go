package middleware

import (
	"forum/api/repository"
	"forum/api/repository/crud"
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

func checkUserRole(id int64) (int, int, error) {
	var (
		role   int
		um     repository.UserRepo
		status int
		err    error
	)
	um = crud.NewUserRepoCRUD()
	if role, status, err = um.GetRole(id); err != nil {
		return role, status, err
	}
	return role, status, nil
}
