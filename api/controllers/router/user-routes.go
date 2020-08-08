package router

import (
	"forum/api/controllers"
	"net/http"
)

var userRoutes = []Route{
	Route{"/users/", http.MethodGet, controllers.UserHandler},
}
