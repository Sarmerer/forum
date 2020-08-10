package controllers

import (
	"forum/api/errors"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		errors.HTTPErrorsHandler(http.StatusNotFound, w, r)
	}
}
