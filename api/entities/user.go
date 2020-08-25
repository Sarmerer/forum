package entities

import "time"

//User struct contains info about user
type User struct {
	ID         uint64
	Name       string
	Password   string
	Email      string
	Nickname   string
	Created    time.Time
	LastOnline time.Time
	SessionID  string
	Role       int
}
