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
		base := []middleware.Middlewares{
			middleware.Logger,
			middleware.SetHeaders,
			middleware.CheckAPIKey,
		}
		if route.NeedAuth {
			base = append(base, middleware.CheckUserAuth)
		}
		if route.SelfOnly {
			base = append(base, middleware.SelfActionOnly)
		}
		mux.HandleFunc(route.URI, route.Method, middleware.Chain(route.Handler, base...))
	}
}
