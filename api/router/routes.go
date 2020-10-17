package router

import (
	"forum/api/config"
	"forum/api/controllers"
	"forum/api/middleware"
	"net/http"
)

type route struct {
	URI      string
	Handler  func(http.ResponseWriter, *http.Request)
	Method   string
	MinRole  int
	NeedAuth bool
}

// SetupRoutes sets handlers with middleware chains to API routes
func (mux *Router) SetupRoutes() {
	routes := apiRoutes
	for _, route := range routes {
		seq := []middleware.Middlewares{
			middleware.Logger,
			middleware.SetHeaders,
			middleware.SetContext,
		}
		if route.NeedAuth {
			seq = append(seq, middleware.AuthorizedOnly)
		}
		if route.MinRole == config.RoleModer {
			seq = append(seq, middleware.ModerOrAdmin)
		}
		if route.MinRole == config.RoleAdmin {
			seq = append(seq, middleware.AdminOnly)
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
		Handler:  controllers.LogIn,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	{
		URI:      "/api/auth/signup",
		Handler:  controllers.SignUp,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	{
		URI:      "/api/auth/signout",
		Handler:  controllers.LogOut,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: true,
	},
	{
		URI:      "/api/auth/me",
		Handler:  controllers.Me,
		Method:   http.MethodGet,
		MinRole:  config.RoleUser,
		NeedAuth: true,
	},

	/* -------------------------------------------------------------------------- */
	/*                                 User routes                                */
	/* -------------------------------------------------------------------------- */

	{
		URI:      "/api/user",
		Handler:  controllers.GetUser,
		Method:   http.MethodGet,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	{
		URI:      "/api/users",
		Handler:  controllers.GetUsers,
		Method:   http.MethodGet,
		MinRole:  config.RoleAdmin,
		NeedAuth: true,
	},
	{
		URI:      "/api/user/update",
		Handler:  controllers.UpdateUser,
		Method:   http.MethodPut,
		MinRole:  config.RoleUser,
		NeedAuth: true,
	},
	{
		URI:      "/api/user/delete",
		Handler:  controllers.DeleteUser,
		Method:   http.MethodDelete,
		MinRole:  config.RoleUser,
		NeedAuth: true,
	},

	/* -------------------------------------------------------------------------- */
	/*                                 Post routes                                */
	/* -------------------------------------------------------------------------- */

	{
		URI:      "/api/post/find",
		Handler:  controllers.FindPost,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	{
		URI:      "/api/posts",
		Handler:  controllers.GetPosts,
		Method:   http.MethodGet,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	{
		URI:      "/api/post/create",
		Handler:  controllers.CreatePost,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: true,
	},
	{
		URI:      "/api/post/update",
		Handler:  controllers.UpdatePost,
		Method:   http.MethodPut,
		MinRole:  config.RoleUser,
		NeedAuth: true,
	},
	{
		URI:      "/api/post/delete",
		Handler:  controllers.DeletePost,
		Method:   http.MethodDelete,
		MinRole:  config.RoleUser,
		NeedAuth: true,
	},
	{
		URI:      "/api/post/rate",
		Handler:  controllers.RatePost,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: true,
	},
	/* -------------------------------------------------------------------------- */
	/*                               Categories routes                            */
	/* -------------------------------------------------------------------------- */
	{
		URI:      "/api/categories",
		Handler:  controllers.GetAllCategories,
		Method:   http.MethodGet,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	/* -------------------------------------------------------------------------- */
	/*                               Comment routes                               */
	/* -------------------------------------------------------------------------- */

	{
		URI:      "/api/comments",
		Handler:  controllers.GetComments,
		Method:   http.MethodGet,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	{
		URI:      "/api/comment/add",
		Handler:  controllers.CreateComment,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: true,
	},
	{
		URI:      "/api/comment/update",
		Handler:  controllers.UpdateComment,
		Method:   http.MethodPut,
		MinRole:  config.RoleUser,
		NeedAuth: true,
	},
	{
		URI:      "/api/comment/delete",
		Handler:  controllers.DeleteComment,
		Method:   http.MethodDelete,
		MinRole:  config.RoleUser,
		NeedAuth: true,
	},
}
