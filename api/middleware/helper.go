package middleware

import (
	"database/sql"
	models "forum/api/models/user"
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
	var role int
	var err error
	var status int
	var db *sql.DB
	var um *models.UserModel
	status = http.StatusInternalServerError
	db, err = database.Connect()
	if err != nil {
		return 0, status, err
	}
	defer db.Close()
	um, err = models.NewUserModel(db)
	if err != nil {
		return 0, status, err
	}
	if role, status, err = um.GetRole(id); err != nil {
		return role, status, err
	}
	return role, status, nil
}
