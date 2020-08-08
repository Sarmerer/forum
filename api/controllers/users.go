package controllers

import (
	"net/http"
	"strings"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUser(w, r)
	case http.MethodPost:
		createUser(w, r)
	case http.MethodPut:
		updateUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get user " + strings.TrimPrefix(r.URL.Path, "/users/")))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create user"))
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update user"))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete user"))
}
