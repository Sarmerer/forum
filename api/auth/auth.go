package auth

import (
	"forum/utils"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/signin":
			signIn(w, r)
		case "/signup":
			signUp(w, r)
		default:
			utils.HTTPErrorsHandler(404, w, r)
		}
	default:
		utils.HTTPErrorsHandler(405, w, r)
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.FormValue("email") + r.FormValue("password")))
}

func signUp(w http.ResponseWriter, r *http.Request) {

}
