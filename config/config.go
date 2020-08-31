package config

import (
	"time"
)

var ()

const (
	DbDriver   = "sqlite3"
	DbPath     = "./database/forum.db"
	APIPort    = 4433
	TimeLayout = "2006-01-02 15:04:05"
	GCInterval = 24 * time.Hour

	//---> Cookies
	SessionCookieName = "session"
	SessionExpiration = 14 * 24 * time.Hour

	//---> Statuses
	StatusSuccess = "success"
	StatusError   = "error"

	//---> Roles
	RoleDefault   = 0
	RoleModerator = 1
	RoleAdmin     = 2
)
