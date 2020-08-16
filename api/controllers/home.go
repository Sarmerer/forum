package controllers

import (
	"net/http"
	"text/template"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("./ui/templates/*html"))
	data := struct {
		Title string
		Next  string
	}{"Ban", "test"}
	tmpl.ExecuteTemplate(w, "base.html", data)
	//w.Write([]byte("get home"))
}
