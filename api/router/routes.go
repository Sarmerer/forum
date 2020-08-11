package router

import (
	"forum/api/auth"
	"forum/api/controllers"
	"net/http"
)

var routes = []Route{
	{Type: "fileServer", URI: "/static/", FileServer: http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static")))},

	{Type: "route", URI: "/signin", Handler: auth.AuthHandler},
	{Type: "route", URI: "/signup", Handler: auth.AuthHandler},
	{Type: "route", URI: "/signout", Handler: auth.AuthHandler},

	{Type: "route", URI: "/", Handler: controllers.RootHandler},
	{Type: "route", URI: "/home", Handler: controllers.HomeHandler},
	{Type: "route", URI: "/users/", Handler: controllers.UsersHandler},
	{Type: "route", URI: "/posts/", Handler: controllers.PostsHandler},
}
