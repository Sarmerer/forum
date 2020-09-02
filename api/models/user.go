package models

type User struct {
	ID          uint64
	Login       string
	Password    string
	Email       string
	DisplayName string
	Created     string
	LastOnline  string
	SessionID   string
	Role        int
}
