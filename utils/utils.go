package utils

import (
	"html/template"
	"net/http"
)

func ServeTemplate(path, layout string, w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseGlob(path+"*.html")).ExecuteTemplate(w, layout, http.StatusOK)
}
