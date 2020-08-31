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
func NewPostReplyModel(db *sql.DB) *PostReplyModel {
	return &PostReplyModel{db}
}

//FindAll returns all replies for the specified post
func (prm *PostReplyModel) FindAll(postID uint64) ([]models.PostReply, error) {
	var (
		rows    *sql.Rows
		replies []models.PostReply
		err     error
	)
	if rows, err = prm.DB.Query("SELECT * FROM replies WHERE reply_post = ?", postID); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	for rows.Next() {
		var reply models.PostReply
		var date string
		rows.Scan(&reply.ID, &reply.Content, &date, &reply.Post, &reply.By)
		if reply.Date, err = time.Parse(config.TimeLayout, date); err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}
	return replies, nil
}

//FindByID returns a specific reply from the database
func (prm *PostReplyModel) FindByID(rid uint64) (*models.PostReply, int, error) {
	var (
		date  string
		reply models.PostReply
		err   error
	)
	if err = prm.DB.QueryRow("SELECT * FROM replies WHERE reply_id = ?", rid).Scan(&reply.ID, &reply.Content, &date, &reply.Post, &reply.By); err == sql.ErrNoRows {
		return nil, http.StatusBadRequest, errors.New("reply not found")
	} else if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if reply.Date, err = time.Parse(config.TimeLayout, date); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &reply, http.StatusOK, nil
}

//Create adds a new reply to the database
func (prm *PostReplyModel) Create(reply *models.PostReply) error {
	var (
		stmt         *sql.Stmt
		res          sql.Result
		rowsAffected int64
		err          error
	)
	if stmt, err = prm.DB.Prepare("INSERT INTO replies (reply_content, reply_date, reply_post, reply_by) VALUES (?, ?, ?, ?)"); err != nil {
		return err
	}
	if res, err = stmt.Exec(reply.Content, time.Now().Format(config.TimeLayout), reply.Post, reply.By); err != nil {
		return err
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected > 0 {
		return nil
	}
	return errors.New("could not create the reply")
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
	if res, err = stmt.Exec(reply.Content, reply.Date.Format(config.TimeLayout), reply.Post, reply.By, reply.ID); err != nil {
		return err
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected > 0 {
		return nil
	}
	return errors.New("failed to update the reply")
}

//Delete deletes reply from the database
func (prm *PostReplyModel) Delete(rid uint64) error {
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
	} else if rowsAffected > 0 {
		return nil
	}
	return errors.New("failed to delete the reply")
}

func (prm *PostReplyModel) DeleteGroup(pid uint64) error {
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
	} else if rowsAffected > 0 {
		return nil
	}
	return errors.New("failed to delete replies for the post")
}
