package middleware

import (
	"net/http"
)

type Middlewares func(http.HandlerFunc) http.HandlerFunc

//Chain function takes in multiple middleware functions,
//and combines them, to prevent spaghetti code
func Chain(h http.HandlerFunc, m ...Middlewares) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}
	wrapped := h
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped
}

//getUserIPAddress : find user IP address from the request
func getUserIPAddress(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
