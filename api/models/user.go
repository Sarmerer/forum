package models

type User struct {
	ID          int64 `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	DisplayName string `json:"display_name"`
	Created     string `json:"created"`
	LastOnline  string `json:"last_online"`
	SessionID   string `json:"session_id"`
	Role        int    `json:"role"`
}
