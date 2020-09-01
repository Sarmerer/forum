package config

import (
	"time"
)

var ()

const (
	// DatabasePath defines the path to the database
	DatabasePath = "./database/forum.db"
	// APIPort defines the port on which the API will run
	APIPort = 4433
	// TimeLayout defines the standartd time layout
	// All time variables are formatted corresponding to it
	TimeLayout = "2006-01-02 15:04:05"
	// GCInterval defines the interval after which garbage collector will run
	GCInterval = 24 * time.Hour

	// SessionCookieName defines the name of the session cookie,
	// which will be stored in client's cookie-jar
	SessionCookieName = "session"
	// SessionExpiration defines the cookie life time
	SessionExpiration = 14 * 24 * time.Hour

	// RoleDefault defines the permission level of a user
	RoleDefault = 0
	// RoleModerator defines the permission level of a user
	RoleModerator = 1
	// RoleAdmin defines the permission level of a user
	RoleAdmin = 2
)
