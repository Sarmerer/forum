package router

import (
	"forum/api/controllers"
	"forum/api/middleware"
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

// SetupRoutes sets handlers with middleware chains to API routes
func (mux *Router) SetupRoutes() {
	routes := apiRoutes
	for _, route := range routes {
		seq := []middleware.Middlewares{
			middleware.Logger,
			middleware.SetHeaders,
			middleware.CheckAPIKey,
		}
		if route.NeedAuth {
			seq = append(seq, middleware.CheckUserAuth)
		}
		if route.SelfOnly {
			seq = append(seq, middleware.SelfActionOnly)
		}
		mux.HandleFunc(route.URI, route.Method, middleware.Chain(route.Handler, seq...))
	}
}

var apiRoutes = []route{

	/* -------------------------------------------------------------------------- */
	/*                                 Auth routes                                */
	/* -------------------------------------------------------------------------- */

	{
		URI:      "/api/auth/signin",
		Handler:  controllers.SignIn,
		Method:   http.MethodPost,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: false,
	},
	{
		URI:      "/api/auth/signup",
		Handler:  controllers.SignUp,
		Method:   http.MethodPost,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: false,
	},
	{
		URI:      "/api/auth/signout",
		Handler:  controllers.SignOut,
		Method:   http.MethodPost,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: true,
	},
	{
		URI:      "/api/auth/status",
		Handler:  controllers.Status,
		Method:   http.MethodPost,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: false,
	},

	/* -------------------------------------------------------------------------- */
	/*                                 User routes                                */
	/* -------------------------------------------------------------------------- */

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

	/* -------------------------------------------------------------------------- */
	/*                                 Post routes                                */
	/* -------------------------------------------------------------------------- */

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

	/* -------------------------------------------------------------------------- */
	/*                               Comment routes                               */
	/* -------------------------------------------------------------------------- */

	{
		URI:      "/api/comment/add",
		Handler:  controllers.CreateReply,
		Method:   http.MethodPost,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: true,
	},
	{
		URI:      "/api/comment/update",
		Handler:  controllers.UpdateReply,
		Method:   http.MethodPut,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: true,
	},
	{
		URI:      "/api/comment/delete",
		Handler:  controllers.DeleteReply,
		Method:   http.MethodDelete,
		MinRole:  0,
		SelfOnly: false,
		NeedAuth: true,
	},
}
