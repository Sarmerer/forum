package router

import "net/http"

type Route struct {
	URI     string
	Handler func(http.ResponseWriter, *http.Request)
}

func New() *http.ServeMux {
	mux := http.NewServeMux()
	return SetupRoutes(mux)
}

func SetupRoutes(mux *http.ServeMux) *http.ServeMux {
	for _, route := range userRoutes {
		mux.HandleFunc(route.URI, route.Handler)
	}
	return mux
}
