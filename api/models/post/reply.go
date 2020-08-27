package models

import (
	"database/sql"
	"errors"
	"forum/config"
	"net/http"
	"time"
)

//PostReply struct contains info about a reply to a post
type PostReply struct {
	ID      int
	Content string
	Date    time.Time
	Post    int64
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
func (pr *PostReplyModel) FindAll(postID int64) ([]PostReply, error) {
	rows, err := pr.DB.Query("SELECT * FROM replies WHERE reply_post = ?", postID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
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
func (pr *PostReplyModel) Find(id int) (PostReply, error) {
	var reply PostReply
	rows, err := pr.DB.Query("SELECT * FROM replies WHERE reply_id = ?", id)
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
func (pr *PostReplyModel) Create(reply *PostReply) (int, error) {
	statement, err := pr.DB.Prepare("INSERT INTO replies (reply_content, reply_date, reply_post, reply_by) VALUES (?, ?, ?, ?)")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	res, err := statement.Exec(reply.Content, reply.Date, reply.Post, reply.By)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return http.StatusOK, nil
	}
	return http.StatusBadRequest, errors.New("could not create the reply")
}

//Delete deletes reply from the database
func (pr *PostReplyModel) Delete(id int) bool {
	res, err := pr.DB.Exec("DELETE FROM replies WHERE reply_id = ?", id)
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
func (pr *PostReplyModel) Update(reply *PostReply) bool {
	statement, err := pr.DB.Prepare("UPDATE replies SET reply_content = ?, reply_date = ?, reply_post = ?, reply_by = ? WHERE reply_id = ?")
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
