package config

import (
	"time"
)

type flag struct {
	Message string
	State   *bool
}

// Actual config variables
var (
	ClientURL = clientURLDev
	APIURL    = apiURLDev
	Flags     = map[string]*flag{"--prod": {"production mode is on", &Production}}
)

// API config
const (
	APIPort       = "4433"
	apiURLDev     = "http://localhost:" + APIPort
	apiURLProd    = "https://forum-api-sarmerer.herokuapp.com"
	clientURLDev  = "http://localhost:8080"
	clientURLProd = "https://forum-sarmerer.herokuapp.com"

	DatabaseDriver   = "sqlite3"
	DatabasePath     = "./database"
	DatabaseFileName = "forum.db"

	GCInterval = 24 * time.Hour // GCInterval defines the interval after which garbage collector will run

	SessionCookieName = "sid"              // SessionCookieName defines the name of the session cookie, which will be stored in client's cookie-jar
	SessionExpiration = 2 * 24 * time.Hour // SessionExpiration defines the session cookie life time

	UserCtxVarName = "userCtx" // Used when setting request context in middleware

	RoleUser  = 0
	RoleModer = 1
	RoleAdmin = 2
)

// Flags
var (
	// When production mode is on, backend api starts to use clientURLProd,
	// instead of clientURLDev, to set CORS header Allow-Origin.
	// It is required for session auth to work properly.
	// Production mode can be activated with --prod flag, when starting the app.
	Production = false
)

func Init() {
	if Production {
		ClientURL = clientURLProd
		APIURL = apiURLProd
	}
}
