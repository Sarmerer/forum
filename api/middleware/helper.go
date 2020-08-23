package middleware

import (
	"net/http"
)

type Middlewares func(http.HandlerFunc) http.HandlerFunc

//Chain function takes in multiple middleware functions,
//and combines them, to avoid spaghetti code.
func Chain(h http.HandlerFunc, m ...Middlewares) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}
	wrapped := h
	for i := 0; i <= len(m)-1; i++ {
		wrapped = m[i](wrapped)
	}
	return wrapped
}
