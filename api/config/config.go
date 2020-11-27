package config

import (
	"time"
)

type Flag struct {
	Message string
	State   *bool
}

// Actual config variables
var (
	ClientURL = clientURLTest
	APIURL    = apiURLTest
	Flags     = map[string]*Flag{"--prod": {"production mode is on", &Production}}
)

// API config
const (
	APIPort       = "4433"
	apiURLTest    = "http://localhost:" + APIPort
	apiURLProd    = "https://forum-api-sarmerer.herokuapp.com"
	clientURLTest = "http://localhost:8080"
	clientURLProd = "https://forum-sarmerer.netlify.app"

	DatabaseDriver   = "sqlite3"
	DatabasePath     = "./database"
	DatabaseFileName = "forum.db"

	TimeLayout = "2006-01-02 15:04:05" // TimeLayout defines the standartd time layout. All time variables are formatted corresponding to it
	GCInterval = 24 * time.Hour        // GCInterval defines the interval after which garbage collector will run

	SessionCookieName = "sid"               // SessionCookieName defines the name of the session cookie, which will be stored in client's cookie-jar
	SessionExpiration = 14 * 24 * time.Hour // SessionExpiration defines the cookie life time

	UserCtxVarName = "userCtx" // Used when setting request context in middleware

	RoleUser  = 0
	RoleModer = 1
	RoleAdmin = 2
)

// Flags
var (
	Production = false
)

func Init() {
	if Production {
		ClientURL = clientURLProd
		APIURL = apiURLProd
	}
}
