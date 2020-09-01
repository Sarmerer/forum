package router

import (
	"errors"
	"forum/api/logger"
	"forum/api/response"
	"net/http"
)

// Router contains a map of API routes, with their handlrers
type Router struct {
	handlers map[string]*handler
}

type handler struct {
	Handler func(http.ResponseWriter, *http.Request)
	Method  string
}

// New creates an instance of Router
func New() *Router {
	router := new(Router)
	router.handlers = make(map[string]*handler)
	return router
}

// ServeHTTP is called for every request, it finds an API route, matching request path, and calls the handler for that path
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

// HandleFunc adds a route pattern to the router
func (r *Router) HandleFunc(path, method string, h http.HandlerFunc) {
	r.handlers[path] = &handler{h, method}
}
