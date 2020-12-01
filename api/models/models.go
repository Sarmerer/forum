package models

type UserCtx struct {
	ID   int64
	Role int
}

type User struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password,omitempty"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	DisplayName string `json:"display_name"`
	Created     string `json:"created,omitempty"`
	LastActive  string `json:"last_active,omitempty"`
	SessionID   string `json:"session_id,omitempty"`
	Role        int    `json:"role"`

	Rating   int `json:"rating"`
	Posts    int `json:"posts"`
	Comments int `json:"comments"`
}

type Post struct {
	ID           int64      `json:"id"`
	AuthorID     int64      `json:"author_id"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	Created      string     `json:"created"`
	Updated      string     `json:"updated"`
	Categories   []Category `json:"categories"`
	Rating       int        `json:"rating"`
	YourReaction int        `json:"your_reaction"`

	Author *User `json:"author"`

	CommentsCount     int `json:"comments_count"`
	ParticipantsCount int `json:"participants_count"`
}

type Posts struct {
	Hot       []Post `json:"hot"`
	Recent    []Post `json:"recent"`
	TotalRows int    `json:"total_rows"`
}

type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UseCount int64  `json:"use_count,omitempty"`
}

type Comment struct {
	ID           int64  `json:"id"`
	AuthorID     int64  `json:"author_id"`
	Content      string `json:"content"`
	Created      string `json:"created"`
	PostID       int64  `json:"post"`
	Rating       int    `json:"rating"`
	YourReaction int    `json:"your_reaction"`
	Edited       string `json:"edited"`

	Author *User `json:"author"`
}

type Rating struct {
	Rating       int `json:"rating"`
	YourReaction int `json:"your_reaction"`
}
