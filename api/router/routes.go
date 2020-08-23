package router

import (
	"forum/api/middleware"
	"forum/api/router/mux"
	"net/http"
)

type MyMux http.ServeMux
type Route struct {
	URI         string
	Handler     func(http.ResponseWriter, *http.Request)
	Method      string
	RequresAuth bool
}

var base = []middleware.Middlewares{
	middleware.CheckAPIKey,
	middleware.SetHeaders,
	middleware.Logger,
}

var baseWithAuth = []middleware.Middlewares{
	middleware.CheckUserAuth,
	middleware.CheckAPIKey,
	middleware.SetHeaders,
	middleware.Logger,
}

func New() *mux.Router {
	mux := mux.NewRouter()
	setupRoutes(mux)
	return mux
}

func setupRoutes(mux *mux.Router) {
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
