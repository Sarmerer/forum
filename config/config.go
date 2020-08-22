package config

import (
	"time"
)

var ()

const (
	DbDriver   = "sqlite3"
	DbURL      = "./database/forum.db"
	APIPort  = 4433
	TimeLayout = "2006-01-02 15:04:05"

	//---> Cookies
	SessionCookieName = "sessionID"
	CookieExpiration  = 14 * 24 * time.Hour

	//---> Statuses
	StatusSuccess = "success"
	StatusError   = "error"
)
