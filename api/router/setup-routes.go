package router

import (
	"forum/api/middleware"
	"forum/api/router/mux"
)

func New() *mux.Router {
	mux := mux.NewRouter()
	setupRoutes(mux)
	return mux
}

func setupRoutes(mux *mux.Router) {
	routes := apiRoutes
	for _, route := range routes {
		sequence := []middleware.Middlewares{
			middleware.Logger,
			middleware.SetHeaders,
			middleware.CheckAPIKey,
		}
		if route.NeedAuth {
			sequence = append(sequence, middleware.CheckUserAuth)
		}
		if route.SelfOnly {
			sequence = append(sequence, middleware.SelfActionOnly)
		}
		mux.HandleFunc(route.URI, route.Method, middleware.Chain(route.Handler, sequence...))
	}
}
