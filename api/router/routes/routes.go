package routes

import (
	"forum/api/middleware"
	"net/http"
)

type Route struct {
	URI         string
	Handler     func(http.ResponseWriter, *http.Request)
	Method      string
	RequresAuth bool
}

func New() *http.ServeMux {
	mux := http.NewServeMux()
	setupFileServers(mux)
	setupRoutes(mux)
	return mux
}

func setupRoutes(mux *http.ServeMux) {
	routes := authRoutes
	routes = append(routes, userRoutes...)
	routes = append(routes, postRoutes...)
	for _, route := range routes {
		mux.HandleFunc(route.URI, middleware.AllowedMethods(route.Method, route.Handler))
	}
}

func setupFileServers(mux *http.ServeMux) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))
}
