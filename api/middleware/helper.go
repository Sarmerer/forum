package middleware

import (
	"forum/api/database"
	"forum/api/models"
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

func checkUserRole(id uint64) (bool, error) {
	db, dbErr := database.Connect()
	defer db.Close()
	if dbErr != nil {
		return false, dbErr
	}
	um, umErr := models.NewUserModel(db)
	if umErr != nil {
		return false, umErr
	}
	if role, grErr := um.GetRole(id); umErr != nil {
		return false, grErr
	} else if role > 0 {
		return true, nil
	}
	return false, nil
}
