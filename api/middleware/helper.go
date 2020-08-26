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

func checkUserRole(id uint64) (role int, err error) {
	var db *sql.DB
	var um *models.UserModel
	db, err = database.Connect()
	defer db.Close()
	if err != nil {
		return
	}
	um, err = models.NewUserModel(db)
	if err != nil {
		return
	}
	if role, err = um.GetRole(id); err != nil {
		return
	}
	return
}
