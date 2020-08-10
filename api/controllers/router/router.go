package router

import "net/http"

type Route struct {
	URI     string
	Handler func(http.ResponseWriter, *http.Request)
}

func New() *http.ServeMux {
	mux := http.NewServeMux()
	setupFileServers(mux)
	setupRoutes(mux)
	return mux
}

func setupRoutes(mux *http.ServeMux) {
	for _, route := range routes {
		mux.HandleFunc(route.URI, route.Handler)
	}
}

func setupFileServers(mux *http.ServeMux) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../ui/static"))))
}
