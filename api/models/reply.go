package models

//PostReply struct contains info about a reply to a post
type PostReply struct {
	ID         uint64 `json:"id"`
	AuthorID   uint64 `json:"author_id"`
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
	Created    string `json:"created"`
	Post       uint64 `json:"post"`
}
