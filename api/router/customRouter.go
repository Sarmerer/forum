package router

import (
	"errors"
	"forum/api/response"
	"net/http"
)

// Router serves http
type Router struct {
	handlers map[string]*handler
}

type handler struct {
	Handler func(http.ResponseWriter, *http.Request)
	Method  string
}

// NewRouter creates instance of Router
func NewRouter() *Router {
	router := new(Router)
	router.handlers = make(map[string]*handler)
	return router
}

// ServeHTTP is called for every connection
func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, found := s.handlers[r.URL.Path]
	if !found {
		response.Error(w, http.StatusNotFound, errors.New("page not found"))
		return
	}
	if f.Method != r.Method {
		response.Error(w, http.StatusMethodNotAllowed, errors.New("wrong method"))
		return
	}
	f.Handler(w, r)
}

func (s *Router) HandleFunc(path, method string, h http.HandlerFunc) {
	s.handlers[path] = &handler{h, method}
}
