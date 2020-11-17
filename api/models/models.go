package models

type UserCtx struct {
	ID          int64
	DisplayName string
	Role        int
}

type User struct {
	ID          int64  `json:"id"`
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

type Post struct {
	ID           int64      `json:"id"`
	AuthorID     int64      `json:"author_id"`
	AuthorName   string     `json:"author_name"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	Created      string     `json:"created"`
	Updated      string     `json:"updated"`
	Categories   []Category `json:"categories,omitempty"`
	Comments     []Comment  `json:"comments,omitempty"`
	Rating       int        `json:"rating"`
	YourReaction int        `json:"your_reaction"`

	CommentsCount       int    `json:"comments_count"`
	TotalParticipants   int    `json:"total_participants"`
	LastCommentFromID   int64  `json:"last_comment_from_id"`
	LastCommentFromName string `json:"last_comment_from_name"`
	LastCommentDate     string `json:"last_comment_date"`
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
	AuthorName   string `json:"author_name"`
	Content      string `json:"content"`
	Created      string `json:"created"`
	PostID       int64  `json:"post"`
	Rating       int    `json:"rating"`
	YourReaction int    `json:"your_reaction"`
	Edited       bool   `json:"edited"`
}

type Rating struct {
	Rating       int `json:"rating"`
	YourReaction int `json:"your_reaction"`
}
