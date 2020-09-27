package models

//Post struct contains info about post
type Post struct {
	ID         uint64
	Author     uint64
	AuthorName string
	Title      string
	Content    string
	Created    string
	Updated    string
	Rating     int
}
