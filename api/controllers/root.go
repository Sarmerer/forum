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

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		response.Error(w, http.StatusNotFound, errors.New("Not found"))
	}
}

func RedirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + fmt.Sprint(strings.Split(r.Host, ":")[0], ":", config.HTTPSport) + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, r, target,
		// see comments below and consider the codes 308, 302, or 301
		http.StatusTemporaryRedirect)
}
