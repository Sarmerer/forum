package crud

import (
	"database/sql"
	"errors"
	"forum/api/models"
	"forum/config"
	"net/http"
	"time"
)

//PostReplyModel helps performing CRUD operations
type PostReplyModel struct {
	DB *sql.DB
}

//NewPostReplyModel creates an instance of PostReplyModel
func NewPostReplyModel(db *sql.DB) (*PostReplyModel, error) {
	return &PostReplyModel{db}, nil
}

//FindAll returns all replies for the specified post
func (prm *PostReplyModel) FindAll(postID uint64) ([]models.PostReply, error) {
	rows, err := prm.DB.Query("SELECT * FROM replies WHERE reply_post = ?", postID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var replies []models.PostReply

	for rows.Next() {
		var reply models.PostReply
		var replyDate string
		rows.Scan(&reply.ID, &reply.Content, &replyDate, &reply.Post, &reply.By)
		date, _ := time.Parse(config.TimeLayout, replyDate)
		reply.Date = date
		replies = append(replies, reply)
	}
	return replies, nil
}

//FindByID returns a specific reply from the database
func (prm *PostReplyModel) FindByID(id int) (models.PostReply, error) {
	var reply models.PostReply
	rows, err := prm.DB.Query("SELECT * FROM replies WHERE reply_id = ?", id)
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
func (prm *PostReplyModel) Create(reply *models.PostReply) (int, error) {
	statement, err := prm.DB.Prepare("INSERT INTO replies (reply_content, reply_date, reply_post, reply_by) VALUES (?, ?, ?, ?)")
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

//Update updates existing reply in the database
func (prm *PostReplyModel) Update(reply *models.PostReply) error {
	var (
		stmt         *sql.Stmt
		res          sql.Result
		rowsAffected int64
		err          error
	)
	if stmt, err = prm.DB.Prepare("UPDATE replies SET reply_content = ?, reply_date = ?, reply_post = ?, reply_by = ? WHERE reply_id = ?"); err != nil {
		return err
	}
	if res, err = stmt.Exec(reply.Content, time.Now().Format(config.TimeLayout), reply.Post, reply.By, reply.ID); err != nil {
		return err
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("failed to update the reply")
}

//Delete deletes reply from the database
func (prm *PostReplyModel) Delete(rid int) error {
	var (
		res          sql.Result
		rowsAffected int64
		err          error
	)
	if res, err = prm.DB.Exec("DELETE FROM replies WHERE reply_id = ?", rid); err != nil {
		return err
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("failed to delete the reply")
}

func (prm *PostReplyModel) DeleteAll(pid uint64) error {
	var (
		res          sql.Result
		rowsAffected int64
		err          error
	)
	if res, err = prm.DB.Exec("DELETE FROM replies WHERE reply_post = ?", pid); err != nil {
		return err
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("failed to delete replies for the post")
}
