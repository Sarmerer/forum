package router

import (
	"fmt"
	"net/http"
	"strings"
)

type Route struct {
	handler func(http.ResponseWriter, *http.Request)
	method  string
	name    string
}

// Router serves http
type Router struct {
	routes []*Route
}

// NewRouter creates instance of Router
func NewRouter() *Router {
	router := new(Router)
	//router.handlers = make(map[string]func(http.ResponseWriter, *http.Request))
	return router
}

// ServeHTTP is called for every connection
func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range s.routes {
		if strings.HasPrefix(r.URL.Path, route.name) && r.Method == route.method {
			route.handler(w, r)
		}
	}
	bad(w)
}

func (r *Router) HandlerFunc(path, method string, f http.HandlerFunc) {
	r.routes = append(r.routes, &Route{f, method, path})
}

// // GET sets get handler
// func (s *Router) GET(path string, f http.HandlerFunc) {
// 	s.handlers[key("GET", path)] = f
// }

// // POST sets post handler
// func (s *Router) POST(path string, f http.HandlerFunc) {
// 	s.handlers[key("POST", path)] = f
// }

// // DELETE sets delete handler
// func (s *Router) DELETE(path string, f http.HandlerFunc) {
// 	s.handlers[key("DELETE", path)] = f
// }

// // PUT sets put handler
// func (s *Router) PUT(path string, f http.HandlerFunc) {
// 	s.handlers[key("PUT", path)] = f
// }

func key(method, path string) string {
	return fmt.Sprintf("%s:%s", method, path)
}

func bad(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"error":"not found"}`))
}
