package router

import (
	"forum/api/auth"
	"forum/api/controllers"
)

var routes = []Route{
	{"/signin", auth.AuthHandler},
	{"/signup", auth.AuthHandler},
	{"/signout", auth.AuthHandler},

	{"/", controllers.RootHandler},
	{"/home", controllers.HomeHandler},
	{"/users/", controllers.UsersHandler},
	{"/posts/", controllers.PostsHandler},
}
