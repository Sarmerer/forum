package controllers

import "net/http"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHome(w, r)
	}
}

func getHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get home"))
}
