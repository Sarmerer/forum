package mux

import (
	"errors"
	"forum/api/logger"
	"forum/api/response"
	"forum/api/utils"
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
		logger.HTTPLogs(utils.PaintStatus(http.StatusNotFound), "0µs", r.Host, utils.PaintMethod(r.Method), r.URL.Path)
		response.Error(w, http.StatusNotFound, errors.New("page not found"))
		return
	}
	if f.Method != r.Method {
		logger.HTTPLogs(utils.PaintStatus(http.StatusMethodNotAllowed), "0µs", r.Host, utils.PaintMethod(r.Method), r.URL.Path)
		response.Error(w, http.StatusMethodNotAllowed, errors.New("wrong method"))
		return
	}
	f.Handler(w, r)
}

func (s *Router) HandleFunc(path, method string, h http.HandlerFunc) {
	s.handlers[path] = &handler{h, method}
}
