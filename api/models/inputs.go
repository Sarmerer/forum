package models

type InputID struct {
	ID int64 `json:"id"`
}

type InputUserSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type InputUserSignUp struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InputPostCreateUpdate struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
}

type InputCommentCreateUpdate struct {
	ID      int64  `json:"id"` // post ID, to which comment is added
	Content string `json:"content"`
}

type InputAllPosts struct {
	PerPage     int    `json:"per_page"`
	CurrentPage int    `json:"current_page"`
	OrderBy     string `json:"order_by"`
	Ascending   bool   `json:"ascending"`
}

type InputFind struct {
	By         string   `json:"by"`
	ID         int64    `json:"id"`
	AuthorID   int64    `json:"author"`
	Categories []string `json:"categories"`
}

type InputRate struct {
	ID       int64 `json:"id"`
	Reaction int   `json:"reaction"`
}
