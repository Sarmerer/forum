package models

type User struct {
	ID         uint64
	Name       string
	Password   string
	Email      string
	Created    string
	LastOnline string
	SessionID  string
	Role       int
}
