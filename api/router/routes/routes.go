package routes

import (
	"forum/api/router"
	"net/http"
)

type Route struct {
	URI         string
	Handler     func(http.ResponseWriter, *http.Request)
	Method      string
	RequresAuth bool
}

func New() *router.Router {
	mux := router.NewRouter()
	//setupFileServers(mux)
	setupRoutes(mux)
	return mux
}

func setupRoutes(mux *router.Router) {
	routes := authRoutes
	routes = append(routes, userRoutes...)
	routes = append(routes, postRoutes...)
	for _, route := range routes {
		mux.HandlerFunc(route.URI, route.Method, route.Handler)
		//mux.HandleFunc(route.URI, middleware.AllowedMethods(route.Handler, route.Method))
	}
}

// func setupFileServers(mux *router.Router) {
// 	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))
// }
