package routes

import (
	"forum/api/controllers"
	"net/http"
)

var postRoutes = []Route{
	{"/posts", controllers.GetPosts, http.MethodGet, false},
	{"/posts/", controllers.GetPost, http.MethodGet, false},
	{"/posts/create", controllers.CreatePost, http.MethodPost, true},
	{"/posts/update/", controllers.UpdatePost, http.MethodPut, true},
	{"/posts/delete/", controllers.DeletePost, http.MethodDelete, true},
}
