package router

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
		if route.RequresAuth {
			mux.HandleFunc(route.URI,
				middleware.Logger(
					middleware.SetJSONType(
						middleware.AllowedMethods(
							route.Method, middleware.CheckUserAuth(route.Handler),
						),
					),
				),
			)
		} else {
			mux.HandleFunc(route.URI,
				middleware.Logger(
					middleware.SetJSONType(
						middleware.AllowedMethods(
							route.Method, route.Handler,
						),
					),
				),
			)
		}
	}
}

func setupFileServers(mux *http.ServeMux) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))
}
