package router

import (
	"forum/api/auth"
	"forum/api/controllers"
	"net/http"
)

var apiRoutes = []Route{
	{"/signin", auth.SignIn, http.MethodPost, false},
	{"/signup", auth.SignUp, http.MethodPost, false},
	{"/signout", auth.SignOut, http.MethodPost, true},

	{"/api/user", controllers.GetUser, http.MethodGet, true},
	{"/api/users", controllers.GetUsers, http.MethodGet, false},
	{"/api/user/update", controllers.UpdateUser, http.MethodPut, true},
	{"/api/user/delete", controllers.DeleteUser, http.MethodDelete, true},

	{"/api/post", controllers.GetPosts, http.MethodGet, false},
	{"/api/posts", controllers.GetPost, http.MethodGet, false},
	{"/api/post/create", controllers.CreatePost, http.MethodPost, true},
	{"/api/post/update", controllers.UpdatePost, http.MethodPut, true},
	{"/api/post/delete", controllers.DeletePost, http.MethodDelete, true},
}
