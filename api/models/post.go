package models

//Post struct contains info about post
type Post struct {
	ID           int64  `json:"id"`
	AuthorID     int64  `json:"author_id"`
	AuthorName   string `json:"author_name"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Created      string `json:"created"`
	Updated      string `json:"updated"`
	Rating       int    `json:"rating"`
	YourReaction int    `json:"your_reaction"`
}
