package mux

import (
	"errors"
	"forum/api/logger"
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
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	f, found := r.handlers[req.URL.Path]
	if !found {
		logger.HTTPLogs(logger.PaintStatus(http.StatusNotFound), "0µs", req.Host, logger.PaintMethod(req.Method), req.URL.Path)
		response.Error(w, http.StatusNotFound, errors.New("page not found"))
		return
	}
	if f.Method != req.Method {
		logger.HTTPLogs(logger.PaintStatus(http.StatusMethodNotAllowed), "0µs", req.Host, logger.PaintMethod(req.Method), req.URL.Path)
		response.Error(w, http.StatusMethodNotAllowed, errors.New("wrong method"))
		return
	}
	f.Handler(w, req)
}

func (r *Router) HandleFunc(path, method string, h http.HandlerFunc) {
	r.handlers[path] = &handler{h, method}
}

func (r *Router) Method(method string) {

}
