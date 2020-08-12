package auth

import (
	"fmt"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println(formatRequest(r))
	cookie, err := r.Cookie("sessionID")
	if r.FormValue("password") == "admin" && r.FormValue("login") == "root" {
		if err == http.ErrNoCookie {
			err = nil
			cookie = &http.Cookie{Name: "sessionID", Value: ""}
		}
		cookie = generateCookie()
		http.SetCookie(w, cookie)
	} else {
		w.Write([]byte("Wrong login or password"))
	}
}
func SignUp(w http.ResponseWriter, r *http.Request) {
	//TODO: create new user
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err == http.ErrNoCookie {
		return
	}
	delete(activeSessions, uuid.FromStringOrNil(cookie.Value))
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
