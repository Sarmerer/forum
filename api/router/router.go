package router

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Request struct {
	Vars map[string]string
	http.Request
}

type Route struct {
	handler func(http.ResponseWriter, *http.Request)
	method  string
	pattern string
}

// Router serves http
type Router struct {
	routes []*Route
}

// NewRouter creates instance of Router
func NewRouter() *Router {
	return new(Router)
}

// ServeHTTP is called for every connection
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	matcher := r.Match(req)
	if matcher != nil {
		matcher.handler(w, req)
	}
	bad(w)
}

func (r *Router) HandleFunc(path, method string, f http.HandlerFunc) {
	r.routes = append(r.routes, &Route{f, method, fmt.Sprint("^", path, "$")})
}

func (r *Router) Match(req *http.Request) *Route {
	for _, route := range r.routes {
		if route.Match(req) {
			return route
		}
	}
	return nil
}

func (r *Route) Match(req *http.Request) bool {
	if req.Method != r.method {
		return false
	}
	oIndex := strings.Index(r.pattern, "{")
	eIndex := strings.Index(r.pattern, "}")
	var varNmae string
	if (oIndex >= 0 && eIndex >= 0) && (oIndex < eIndex) {
		varName = r.pattern[oIndex+1 : eIndex]
	}
	pattern := strings.Split(r.pattern, "/")
	request := strings.Split(req.URL.Path, "/")
	if len(pattern) != len(request) {
		return false
	}
	matched, _ := regexp.MatchString(r.pattern, req.URL.Path)
	if matched {
		return true
	}
	return false
}

func (r *Request) Clone(req *http.Request) {

}

func key(method, path string) string {
	return fmt.Sprintf("%s:%s", method, path)
}

func bad(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"error":"not found"}`))
}
