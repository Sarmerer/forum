package models

import (
	"database/sql"
	"forum/config"
	"time"
)

//PostReply struct contains info about a reply to a post
type PostReply struct {
	ID      int
	Content string
	Date    time.Time
	Post    int
	By      int
}

//PostReplyModel helps performing CRUD operations
type PostReplyModel struct {
	DB *sql.DB
}

//NewPostReplyModel creates an instance of PostReplyModel
func NewPostReplyModel(db *sql.DB) (*PostReplyModel, error) {
	return &PostReplyModel{db}, nil
}

//FindAll returns all replies in the database
func (um *PostReplyModel) FindAll() ([]PostReply, error) {
	rows, e := um.DB.Query("SELECT * FROM replies")
	if e != nil {
		return nil, e
	}
	var replies []PostReply

	for rows.Next() {
		var reply PostReply
		var replyDate string
		rows.Scan(&reply.ID, &reply.Content, &replyDate, &reply.Post, &reply.By)
		date, _ := time.Parse(config.TimeLayout, replyDate)
		reply.Date = date
		replies = append(replies, reply)
	}
	return replies, nil
}

//Find returns a specific reply from the database
func (um *PostReplyModel) Find(id int) (PostReply, error) {
	var reply PostReply
	rows, err := um.DB.Query("SELECT * FROM replies WHERE reply_id = ?", id)
	if err != nil {
		return reply, err
	}
	for rows.Next() {
		var replyDate string
		rows.Scan(&reply.ID, &reply.Content, &replyDate, &reply.Post, &reply.By)
		date, _ := time.Parse(config.TimeLayout, replyDate)
		reply.Date = date
	}
	return reply, nil
}

//Create adds a new reply to the database
func (um *PostReplyModel) Create(reply *PostReply) (bool, string) {
	statement, err := um.DB.Prepare("INSERT INTO replies (reply_content, reply_date, reply_post, reply_by) VALUES (?, ?, ?, ?)")
	if err != nil {
		return false, "Internal server error"
	}
	res, err := statement.Exec(reply.Content, time.Now().Format(config.TimeLayout), reply.Post, reply.By)
	if err != nil {
		return false, err.Error()
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, "Internal server error"
	}
	if rowsAffected > 0 {
		return true, ""
	}
	return false, "Internal server error"
}

//Delete deletes reply from the database
func (um *PostReplyModel) Delete(id int) bool {
	res, err := um.DB.Exec("DELETE FROM replies WHERE reply_id = ?", id)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}

//Update updates existing reply in the database
func (um *PostReplyModel) Update(reply *PostReply) bool {
	statement, err := um.DB.Prepare("UPDATE replies SET reply_content = ?, reply_date = ?, reply_post = ?, reply_by = ? WHERE reply_id = ?")
	if err != nil {
		return false
	}
	res, err := statement.Exec(reply.Content, time.Now().Format(config.TimeLayout), reply.Post, reply.By, reply.ID)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}
