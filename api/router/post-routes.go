package router

import (
	"forum/api/controllers"
	"net/http"
)

var postRoutes = []Route{
	{"/post/", controllers.GetPost, http.MethodGet, true},
	{"/posts", controllers.GetPosts, http.MethodGet, false},

	{"/new-post", controllers.CreatePost, http.MethodPost, true},
	{"/update-post", controllers.UpdatePost, http.MethodPut, true},
	{"/delete-post", controllers.DeletePost, http.MethodDelete, true},
}
