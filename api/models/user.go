package models

import "time"

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
