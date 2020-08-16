package controllers

import (
	"net/http"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get home"))
}
