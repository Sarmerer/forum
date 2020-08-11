package router

import (
	"net/http"
)

type Route struct {
	Type       string
	URI        string
	Handler    func(http.ResponseWriter, *http.Request)
	FileServer http.Handler
}

func New() *http.ServeMux {
	mux := http.NewServeMux()
	setupRoutes(mux)
	return mux
}

func setupRoutes(mux *http.ServeMux) {
	for _, route := range routes {
		switch route.Type {	
		case "fileServer":
			mux.Handle("/static/", route.FileServer)
		case "route":
			mux.HandleFunc(route.URI, route.Handler)
		}
	}
}
