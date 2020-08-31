package models

import "time"

//PostReply struct contains info about a reply to a post
type PostReply struct {
	ID      uint64
	Content string
	Date    time.Time
	Post    uint64
	By      uint64
}
