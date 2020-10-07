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
	Title      string
	Content    string
	Categories []string
}

type InputPostFind struct {
	By         string   `json:"by"`
	ID         int64    `json:"id"`
	Author     int64    `json:"author"`
	Categories []string `json:"categories"`
}
