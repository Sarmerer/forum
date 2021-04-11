package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/sarmerer/forum/api/logger"
	"github.com/sarmerer/forum/api/response"
)

// Router contains a map of API routes, with their handlrers
type Router struct {
	routes []*Route
}

type Route struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
	Method  string
}

// New creates an instance of Router
func New() *Router {
	return new(Router)
}

// ServeHTTP is called for every request, it finds an API endpoint, matching request path, and calls the handler for that path
func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range router.routes {
		re := regexp.MustCompile(route.Pattern)
		match := re.MatchString(req.URL.Path)
		if match {
			if req.Method != route.Method && req.Method != http.MethodOptions {
				wrongMethod(w, req)
				return
			}
			route.Handler(w, req)
			return
		}
	}
	pageNotFound(w, req)
	return
}

// HandleFunc adds a route pattern to the router
func (router *Router) HandleFunc(path, method string, handler http.HandlerFunc) {
	var key string = fmt.Sprintf("^(%s)$", path)
	if router.pathExists(key, method) {
		log.Fatalf("duplicated routes: path %s method %s\n", path, method)
		os.Exit(1)
	}
	router.routes = append(router.routes, &Route{key, handler, method})
}

func (router *Router) Handle(path, method string, h http.Handler) {
	var key string = fmt.Sprintf("^(%s)$", path)
	if router.pathExists(key, method) {
		log.Fatalf("duplicated routes: path %s method %s\n", path, method)
		os.Exit(1)
	}
	router.routes = append(router.routes, &Route{key, func(w http.ResponseWriter, r *http.Request) { h.ServeHTTP(w, r) }, method})
}

func pageNotFound(w http.ResponseWriter, req *http.Request) {
	logger.HTTPLogs(logger.PaintStatus(http.StatusNotFound), "0µs", req.Host, logger.PaintMethod(req.Method), req.URL.Path)
	response.Error(w, http.StatusNotFound, errors.New("page not found"))
}

func wrongMethod(w http.ResponseWriter, req *http.Request) {
	logger.HTTPLogs(logger.PaintStatus(http.StatusMethodNotAllowed), "0µs", req.Host, logger.PaintMethod(req.Method), req.URL.Path)
	response.Error(w, http.StatusMethodNotAllowed, errors.New("wrong method"))
}

func (router *Router) pathExists(pattern, method string) bool {
	for _, r := range router.routes {
		if r.Pattern == pattern && r.Method == method {
			return true
		}
	}
	return false
}
