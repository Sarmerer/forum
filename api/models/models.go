package models

type UserCtx struct {
	ID   int64
	Role int
}

type User struct {
	ID          int64  `json:"id" bson:"id,omitempty"`
	Login       string `json:"login" bson:"login,omitempty"`
	Password    string `json:"password,omitempty" bson:"password,omitempty"`
	Email       string `json:"email" bson:"email,omitempty"`
	Avatar      string `json:"avatar" bson:"avatar,omitempty"`
	DisplayName string `json:"display_name" bson:"display_name,omitempty"`
	Created     int64  `json:"created,omitempty" bson:"created,omitempty"`
	LastActive  int64  `json:"last_active,omitempty" bson:"last_active,omitempty"`
	SessionID   string `json:"session_id,omitempty" bson:"session_id,omitempty"`
	Role        int    `json:"role" bson:"role,omitempty"`

	Rating   int `json:"rating"`
	Posts    int `json:"posts"`
	Comments int `json:"comments"`
}

type Post struct {
	ID           int64      `json:"id"`
	AuthorID     int64      `json:"author_id"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	Created      int64      `json:"created"`
	Updated      int64      `json:"updated"`
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
	Created      int64  `json:"created"`
	PostID       int64  `json:"post"`
	Rating       int    `json:"rating"`
	YourReaction int    `json:"your_reaction"`
	Edited       int64  `json:"edited"`

	Author *User `json:"author"`
}

type Rating struct {
	Rating       int `json:"rating"`
	YourReaction int `json:"your_reaction"`
}
