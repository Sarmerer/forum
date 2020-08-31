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
		seq := []middleware.Middlewares{
			middleware.Logger,
			middleware.SetHeaders,
			middleware.CheckAPIKey,
		}
		if route.NeedAuth {
			seq = append(seq, middleware.CheckUserAuth)
		}
		if route.SelfOnly {
			seq = append(seq, middleware.SelfActionOnly)
		}
		mux.HandleFunc(route.URI, route.Method, middleware.Chain(route.Handler, seq...))
	}
}
