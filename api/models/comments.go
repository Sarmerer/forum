package models

//PostComment struct contains info about a reply to a post
type PostComment struct {
	ID         int64 `json:"id"`
	AuthorID   int64 `json:"author_id"`
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
	Created    string `json:"created"`
	Post       int64 `json:"post"`
	Edited     int    `json:"edited"`
}
