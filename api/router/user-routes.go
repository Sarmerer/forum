package router

import (
	"forum/api/controllers"
	"net/http"
)

var userRoutes = []Route{
	{"/api/user", controllers.GetUser, http.MethodGet, true},
	{"/api/users", controllers.GetUsers, http.MethodGet, false},
	{"/api/users/delete", controllers.DeleteUser, http.MethodDelete, true},
	{"/api/users/update", controllers.UpdateUser, http.MethodPut, true},
}
