package router

import (
	"errors"
	"net/http"

	"github.com/sarmerer/forum/api/logger"
	"github.com/sarmerer/forum/api/response"
)

// Router contains a map of API routes, with their handlrers
type Router struct {
	handlers map[string]*routeHandler
}

type routeHandler struct {
	Handler func(http.ResponseWriter, *http.Request)
	Method  string
}

// New creates an instance of Router
func New() *Router {
	router := new(Router)
	router.handlers = make(map[string]*routeHandler)
	return router
}

// ServeHTTP is called for every request, it finds an API endpoint, matching request path, and calls the handler for that path
func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route, found := router.handlers[req.Method+":"+req.URL.Path]
	if !found {
		pageNotFound(w, req)
		return
	}
	if route.Method != req.Method && req.Method != http.MethodOptions {
		wrongMethod(w, req)
		return
	}
	route.Handler(w, req)
}

// HandleFunc adds a route pattern to the router
func (router *Router) HandleFunc(path, method string, h http.HandlerFunc) {
	router.handlers[method+":"+path] = &routeHandler{h, method}
}

func pageNotFound(w http.ResponseWriter, req *http.Request) {
	logger.HTTPLogs(logger.PaintStatus(http.StatusNotFound), "0µs", req.Host, logger.PaintMethod(req.Method), req.URL.Path)
	response.Error(w, http.StatusNotFound, errors.New("page not found"))
}

func wrongMethod(w http.ResponseWriter, req *http.Request) {
	logger.HTTPLogs(logger.PaintStatus(http.StatusMethodNotAllowed), "0µs", req.Host, logger.PaintMethod(req.Method), req.URL.Path)
	response.Error(w, http.StatusMethodNotAllowed, errors.New("wrong method"))
}
