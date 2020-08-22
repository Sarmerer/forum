package router

import (
	"forum/api/controllers"
	"net/http"
)

var postRoutes = []Route{
	{"/api/posts", controllers.GetPosts, http.MethodGet, false},
	{"/api/posts/", controllers.GetPost, http.MethodGet, false},
	{"/api/posts/create", controllers.CreatePost, http.MethodPost, true},
	{"/api/posts/update/", controllers.UpdatePost, http.MethodPut, true},
	{"/api/posts/delete/", controllers.DeletePost, http.MethodDelete, true},
}
