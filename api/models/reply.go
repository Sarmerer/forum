package models

//PostReply struct contains info about a reply to a post
type PostReply struct {
	ID      uint64
	Author  uint64
	Content string
	Created string
	Post    uint64
}
