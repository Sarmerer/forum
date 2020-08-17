package controllers

import (
	"errors"
	"forum/api/response"
	"net/http"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Write([]byte("get home"))
	} else {
		response.Error(w, http.StatusNotFound, errors.New("not found"))
	}

}
