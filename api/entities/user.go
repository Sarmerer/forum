package entities

import "time"

//User struct contains info about user
type User struct {
	ID         int
	Name       string
	Password   string
	Email      string
	Nickname   string
	Created    time.Time
	LastOnline time.Time
	SessionID  string
	Role       int
}
