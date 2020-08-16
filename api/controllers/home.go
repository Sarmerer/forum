package controllers

import (
	"errors"
	"fmt"
	"forum/api/response"
	"forum/config"
	"log"
	"net/http"
	"strings"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Write([]byte("get home"))
	} else {
		response.Error(w, http.StatusNotFound, errors.New("not found"))
	}

}
