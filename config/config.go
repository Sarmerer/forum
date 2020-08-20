package config

import "time"

var ()

const (
	DbDriver   = "sqlite3"
	DbURL      = "./database/forum.db"
	HTTPSport  = 4433
	HTTPport   = 8080
	TimeLayout = "2006-01-02 15:04:05"

	//#---> Cookies
	SessionCookieName = "sessionID"
	CookieExpiration  = 14 * 24 * time.Hour

	//#---> Statuses
	StatusSuccess = "success"
	StatusError   = "error"
)
