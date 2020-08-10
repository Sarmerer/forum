package controllers

import (
	"forum/utils"
	"html/template"
	"net/http"
)

var tpl *template.Template

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHome(w, r)
	}
}

func getHome(w http.ResponseWriter, r *http.Request) {
	utils.ServeTemplate("../ui/templates/", "login.html", w, r)
	//w.Write([]byte("get home"))
}
