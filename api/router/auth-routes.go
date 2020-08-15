package router

import (
	"forum/api/auth"
	"net/http"
)

var authRoutes = []Route{
	{"/signin", auth.SignIn, http.MethodPost, false},
	{"/signup", auth.SignUp, http.MethodPost, false},
	{"/signout", auth.SignOut, http.MethodPost, true},
}
