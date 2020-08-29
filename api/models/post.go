package models

import "time"

//Post struct contains info about post
type Post struct {
	ID       uint64
	By       uint64
	Category int
	Name     string
	Content  string
	Created  time.Time
	Updated  time.Time
	Rating   int
}
