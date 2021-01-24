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

	AdminToken string `json:"admin_token"` // To verify that user is able to be an admin
}

type InputPostCreateUpdate struct {
	ID         int64    `json:"id"`         //update
	Title      string   `json:"title"`      //create/update
	Content    string   `json:"content"`    //create/update
	Categories []string `json:"categories"` //create/update
}

type InputCommentCreateUpdate struct {
	PostID  int64   `json:"post_id"` // post ID, to which comment is added
	ID      int64   `json:"id"`
	Content string  `json:"content"`
	Parent  Comment `json:"parent"`
}

type InputAllPosts struct {
	PerPage     int    `json:"per_page"`
	CurrentPage int    `json:"current_page"`
	OrderBy     string `json:"order_by"`
	Ascending   bool   `json:"ascending"`
}

type InputFindPost struct {
	By         string   `json:"by"`
	ID         int64    `json:"id"`
	AuthorID   int64    `json:"author"`
	Categories []string `json:"categories"`
}

type InputFindUser struct {
	By string `json:"by"`
	ID int64  `json:"id"`
}

type InputFindComments struct {
	PostID int64 `json:"post_id"`
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
}

type InputRate struct {
	ID       int64 `json:"id"`
	Reaction int   `json:"reaction"`
}
