package router

import (
	"forum/api/controllers"
)

var userRoutes = []Route{
	{"/users/", controllers.UsersHandler},
	{"/home", controllers.HomeHandler},
	{"/posts/", controllers.PostsHandler},
}
