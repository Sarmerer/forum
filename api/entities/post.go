package entities

import "time"

//Post struct contains info about post
type Post struct {
	ID       int64
	By       uint64
	Category int
	Name     string
	Content  string
	Created  time.Time
	Updated  time.Time
	Rating   int
}

//PostReply struct contains info about a reply to a post
type PostReply struct {
	ID      int
	Content string
	Date    time.Time
	Post    int
	By      int
}

//PostReaction struct contains info about a reaction to a post
type PostReaction struct {
	ID       int
	PostID   int
	UserID   int
	Reaction int
}
