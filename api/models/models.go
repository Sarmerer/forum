package models

type UserCtx struct {
	ID   int64
	Role int
}

type User struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	Alias      string `json:"alias"`
	Created    int64  `json:"created,omitempty"`
	LastActive int64  `json:"last_active,omitempty"`
	SessionID  string `json:"session_id,omitempty"`
	Role       int    `json:"role"`

	Rating   int `json:"rating"`
	Posts    int `json:"posts"`
	Comments int `json:"comments"`
}

type Post struct {
	ID           int64      `json:"id"`
	AuthorID     int64      `json:"-"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	Created      int64      `json:"created"`
	Edited       int64      `json:"edited"`
	EditReason   string     `json:"edit_reason"`
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
	PostID       int64  `json:"post_id"`
	AuthorID     int64  `json:"-"`
	ParentID     int64  `json:"parent_id"`
	Depth        int    `json:"depth"`
	Lineage      string `json:"lineage"`
	Content      string `json:"content"`
	Created      int64  `json:"created"`
	Rating       int    `json:"rating"`
	YourReaction int    `json:"your_reaction"`
	Edited       int64  `json:"edited"`
	Deleted      bool   `json:"deleted,omitempty"`

	Author *User `json:"author"`

	Children    []*Comment `json:"children"`
	ChildrenLen int64      `json:"children_length"`
}

type Rating struct {
	Rating       int `json:"rating"`
	YourReaction int `json:"your_reaction"`
}
