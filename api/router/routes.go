package router

import (
	"net/http"
	"time"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/controllers"
	"github.com/sarmerer/forum/api/middleware"
)

type route struct {
	URI      string
	Handler  func(http.ResponseWriter, *http.Request)
	Method   string
	MinRole  int
	NeedAuth bool
	Activity bool

	RateLimit limit
}

type limit struct {
	Requests int
	PerTime  time.Duration
	Cooldown time.Duration
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
		if (limit{} != route.RateLimit) {
			limiter := middleware.RateLimit(route.RateLimit.Requests, route.RateLimit.PerTime, route.RateLimit.Cooldown)
			seq = append(seq, limiter)
		}
		if route.NeedAuth {
			seq = append(seq, middleware.AuthorizedOnly)
		}
		if route.MinRole == config.RoleModer {
			seq = append(seq, middleware.ModerOrHigher)
		}
		if route.MinRole == config.RoleAdmin {
			seq = append(seq, middleware.AdminOnly)
		}
		if route.Activity {
			seq = append(seq, middleware.UpdateUserActivity)
		}
		mux.HandleFunc(route.URI, route.Method, middleware.Chain(route.Handler, seq...))
	}
}

var apiRoutes = []route{

	/* -------------------------------------------------------------------------- */
	/*                                 Auth routes                                */
	/* -------------------------------------------------------------------------- */

	{
		URI:      "/api/oauth",
		Handler:  controllers.OAuthHandler,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	{
		URI:       "/api/auth/verify",
		Handler:   controllers.VerifyEmail,
		Method:    http.MethodPost,
		MinRole:   config.RoleUser,
		NeedAuth:  false,
		RateLimit: limit{Requests: 3, PerTime: time.Second},
	},
	{
		URI:       "/api/auth/send-verification",
		Handler:   controllers.SendVerification,
		Method:    http.MethodPost,
		MinRole:   config.RoleUser,
		NeedAuth:  false,
		RateLimit: limit{Requests: 1, PerTime: time.Minute, Cooldown: 1 * time.Minute},
	},
	{
		URI:       "/api/auth/signin",
		Handler:   controllers.SignIn,
		Method:    http.MethodPost,
		MinRole:   config.RoleUser,
		NeedAuth:  false,
		RateLimit: limit{Requests: 10, PerTime: time.Minute, Cooldown: 2 * time.Minute},
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
		URI:      "/api/users",
		Handler:  controllers.GetUsers,
		Method:   http.MethodGet,
		MinRole:  config.RoleAdmin,
		NeedAuth: true,
	},
	{
		URI:      "/api/user/find",
		Handler:  controllers.FindUser,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	{
		URI:      "/api/user/update",
		Handler:  controllers.UpdateUser,
		Method:   http.MethodPut,
		MinRole:  config.RoleUser,
		NeedAuth: true,
		Activity: true,
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
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	{
		URI:      "/api/post/create",
		Handler:  controllers.CreatePost,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: true,
		Activity: true,
	},
	{
		URI:      "/api/post/update",
		Handler:  controllers.UpdatePost,
		Method:   http.MethodPut,
		MinRole:  config.RoleUser,
		NeedAuth: true,
		Activity: true,
	},
	{
		URI:      "/api/post/delete",
		Handler:  controllers.DeletePost,
		Method:   http.MethodDelete,
		MinRole:  config.RoleUser,
		NeedAuth: true,
		Activity: true,
	},
	{
		URI:      "/api/post/rate",
		Handler:  controllers.RatePost,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: true,
		Activity: true,
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
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	{
		URI:      "/api/comments/find",
		Handler:  controllers.FindComments,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: false,
	},
	{
		URI:      "/api/comment/add",
		Handler:  controllers.CreateComment,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: true,
		Activity: true,
	},
	{
		URI:      "/api/comment/update",
		Handler:  controllers.UpdateComment,
		Method:   http.MethodPut,
		MinRole:  config.RoleUser,
		NeedAuth: true,
		Activity: true,
	},
	{
		URI:      "/api/comment/delete",
		Handler:  controllers.DeleteComment,
		Method:   http.MethodDelete,
		MinRole:  config.RoleUser,
		NeedAuth: true,
		Activity: true,
	},
	{
		URI:      "/api/comment/rate",
		Handler:  controllers.RateComment,
		Method:   http.MethodPost,
		MinRole:  config.RoleUser,
		NeedAuth: true,
		Activity: true,
	},
	/* -------------------------------------------------------------------------- */
	/*                                 Images server                              */
	/* -------------------------------------------------------------------------- */
	{
		URI:     "/api/images",
		Handler: controllers.ServeImage,
		Method:  http.MethodGet,
		MinRole: config.RoleUser,
	},
	{
		URI:     "/api/image/upload",
		Handler: controllers.UploadImage,
		Method:  http.MethodPost,
		MinRole: config.RoleUser,
	},
}
