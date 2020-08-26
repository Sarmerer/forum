package router

import (
	"forum/api/auth"
	"forum/api/controllers"
	"net/http"
)

type route struct {
	URI      string
	Handler  func(http.ResponseWriter, *http.Request)
	Method   string
	MinRole  int
	SelfOnly bool
	NeedAuth bool
}

var apiRoutes = []route{
	//######################################################################
	//############################Auth routes###############################
	//######################################################################
	{
		URI:      "/api/signin",
		Handler:  auth.SignIn,
		Method:   http.MethodPost,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: false,
	},
	{
		URI:      "/api/signup",
		Handler:  auth.SignUp,
		Method:   http.MethodPost,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: false,
	},
	{
		URI:      "/api/signout",
		Handler:  auth.SignOut,
		Method:   http.MethodPost,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: true,
	},
	//######################################################################
	//###########################Users routes###############################
	//######################################################################
	{
		URI:      "/api/user",
		Handler:  controllers.GetUser,
		Method:   http.MethodGet,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: true,
	},
	{
		URI:      "/api/users",
		Handler:  controllers.GetUsers,
		Method:   http.MethodGet,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: false,
	},
	{
		URI:      "/api/user/update",
		Handler:  controllers.UpdateUser,
		Method:   http.MethodPut,
		MinRole:  0,
		SelfOnly: true,
		NeedAuth: true,
	},
	{
		URI:      "/api/user/delete",
		Handler:  controllers.DeleteUser,
		Method:   http.MethodDelete,
		MinRole:  0,
		SelfOnly: true,
		NeedAuth: true,
	},
	//######################################################################
	//###########################Posts routes###############################
	//######################################################################
	{
		URI:      "/api/post",
		Handler:  controllers.GetPost,
		Method:   http.MethodGet,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: false,
	},
	{
		URI:      "/api/posts",
		Handler:  controllers.GetPosts,
		Method:   http.MethodGet,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: false,
	},
	{
		URI:      "/api/post/create",
		Handler:  controllers.CreatePost,
		Method:   http.MethodPost,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: true,
	},
	{
		URI:      "/api/post/update",
		Handler:  controllers.UpdatePost,
		Method:   http.MethodPut,
		MinRole:  0,
		SelfOnly: true,
		NeedAuth: true,
	},
	{
		URI:      "/api/post/delete",
		Handler:  controllers.DeletePost,
		Method:   http.MethodDelete,
		MinRole:  0,
		SelfOnly: true,
		NeedAuth: true,
	},
}
