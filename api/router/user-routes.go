package router

import (
	"forum/api/controllers"
	"net/http"
)

var userRoutes = []Route{
	{"/", controllers.RootHandler, http.MethodGet, false},
	{"/home", controllers.GetHome, http.MethodGet, false},

	{"/users", controllers.GetUsers, http.MethodGet, false},
	{"/users/", controllers.GetUser, http.MethodGet, false},
	{"/users/delete/", controllers.DeleteUser, http.MethodDelete, true},
	{"/users/update/", controllers.UpdateUser, http.MethodPut, true},
}
