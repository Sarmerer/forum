package router

import (
	"forum/api/middleware"
	"net/http"
)

type MyMux http.ServeMux
type Route struct {
	URI         string
	Handler     func(http.ResponseWriter, *http.Request)
	Method      string
	RequresAuth bool
}

var baseWithAuth = []middleware.Middlewares{
	middleware.Logger,
	middleware.SetHeaders,
	middleware.CheckAPIKey,
	middleware.CheckUserAuth,
}

var base = []middleware.Middlewares{
	middleware.Logger,
	middleware.SetHeaders,
	middleware.CheckAPIKey,
}

func New() *Router {
	mux := NewRouter()
	setupRoutes(mux)
	return mux
}

func setupRoutes(mux *Router) {
	routes := authRoutes
	routes = append(routes, userRoutes...)
	routes = append(routes, postRoutes...)
	for _, route := range routes {
		if route.RequresAuth {
			mux.HandleFunc(route.URI, route.Method, middleware.Chain(route.Handler, baseWithAuth...))
		} else {
			mux.HandleFunc(route.URI, route.Method, middleware.Chain(route.Handler, base...))
		}
	}
}
