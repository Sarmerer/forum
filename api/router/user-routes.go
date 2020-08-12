package router

import (
	"forum/api/controllers"
	"net/http"
)

var userRoutes = []Route{
	{"/", controllers.RootHandler, http.MethodGet, false},
	{"/home", controllers.HomeHandler, http.MethodGet, false},

	{"/user/", controllers.GetUser, http.MethodGet, false},
	{"/users", controllers.GetUsers, http.MethodGet, false},

	{"/new-user", controllers.CreateUser, http.MethodGet, false},
	{"/delete-user", controllers.DeleteUser, http.MethodGet, false},
	{"/update-user", controllers.UpdateUser, http.MethodGet, false},
}
