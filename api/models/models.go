package models

type UserCtx struct {
	ID   int64
	Role int
}

type User struct {
	// public data
	ID            int64  `json:"id"`
	Avatar        string `json:"avatar"`
	Alias         string `json:"alias"`
	Created       int64  `json:"created,omitempty"`
	LastActive    int64  `json:"last_active,omitempty"`
	Role          int    `json:"role"`
	OAuthProvider string `json:"oauth_provider"`

	// private data
	// some fields are omitted from json for security reasons
	Username  string `json:"username"`
	Email     string `json:"email"`
	SessionID string `json:"-"`
	Password  string `json:"-"`
	Verified  bool   `json:"-"`

	// stats
	Rating int `json:"rating"`

	Posts          int `json:"posts"`
	PostsUpvoted   int `json:"posts_upvoted"`
	PostsDownvoted int `json:"posts_downvoted"`

	Comments          int `json:"comments"`
	CommentsUpvoted   int `json:"comments_upvoted"`
	CommentsDownvoted int `json:"comments_downvoted"`
}

type Post struct {
	ID       int64  `json:"id"`
	AuthorID int64  `json:"-"`
	Title    string `json:"title"`
	Content  string `json:"content"`

	// Content field is used as an image src, when IsImage is true
	IsImage    bool       `json:"is_image"`
	Created    int64      `json:"created"`
	Edited     int64      `json:"edited"`
	EditReason string     `json:"edit_reason"`
	Categories []Category `json:"categories"`
	Rating     int        `json:"rating"`

	// reaction of a user with x ID, where x is an ID from request context.
	// See UserCtx model
	YourReaction int `json:"your_reaction"`

	Author *User `json:"author"`

	CommentsCount     int `json:"comments_count"`
	ParticipantsCount int `json:"participants_count"`
}

// Posts is a helper type for Post. It allows to include additional data to response
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

// Comments is a helper type for Comment. It allows to include additional data to response
type Comments struct {
	Comments   []*Comment `json:"comments"`
	LoadedRows int        `json:"loaded_rows"`
	TotalRows  int        `json:"total_rows"`
}

type Rating struct {
	Rating       int `json:"rating"`
	YourReaction int `json:"your_reaction"`
}
