package crud

import (
	"database/sql"
	"errors"
	"forum/api/models"
	"forum/config"
	"net/http"
	"time"
)

//ReplyRepoCRUD helps performing CRUD operations
type ReplyRepoCRUD struct {
	DB *sql.DB
}

//NewReplyRepoCRUD creates an instance of PostReplyModel
func NewReplyRepoCRUD(db *sql.DB) *ReplyRepoCRUD {
	return &ReplyRepoCRUD{db}
}

//FindAll returns all replies for the specified post
func (repo *ReplyRepoCRUD) FindAll(postID uint64) ([]models.PostReply, error) {
	var (
		rows    *sql.Rows
		replies []models.PostReply
		err     error
	)
	if rows, err = repo.DB.Query("SELECT * FROM replies WHERE post_id_fkey = ?", postID); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r models.PostReply
		rows.Scan(&r.ID, &r.Author, &r.Content, &r.Created, &r.Post)
		replies = append(replies, r)
	}
	return replies, nil
}

//FindByID returns a specific reply from the database
func (repo *ReplyRepoCRUD) FindByID(rid uint64) (*models.PostReply, int, error) {
	var (
		r   models.PostReply
		err error
	)
	if err = repo.DB.
		QueryRow("SELECT * FROM replies WHERE id = ?", rid).
		Scan(&r.ID, &r.Author, &r.Content, &r.Created, &r.Post); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusBadRequest, errors.New("reply not found")
	}
	return &r, http.StatusOK, nil
}

//Create adds a new reply to the database
func (repo *ReplyRepoCRUD) Create(r *models.PostReply) error {
	var (
		stmt         *sql.Stmt
		res          sql.Result
		rowsAffected int64
		err          error
	)
	if stmt, err = repo.DB.Prepare("INSERT INTO replies (author, content, created, post_id_fkey) VALUES (?, ?, ?, ?)"); err != nil {
		return err
	}
	if res, err = stmt.Exec(r.Author, r.Content, time.Now().Format(config.TimeLayout), r.Post); err != nil {
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
func (repo *ReplyRepoCRUD) Update(r *models.PostReply) error {
	var (
		stmt         *sql.Stmt
		res          sql.Result
		rowsAffected int64
		err          error
	)
	if stmt, err = repo.DB.Prepare("UPDATE replies SET author = ?, content = ?, created = ?, post_id_fkey = ? WHERE reply_id = ?"); err != nil {
		return err
	}
	if res, err = stmt.Exec(r.Author, r.Content, r.Created, r.Post, r.ID); err != nil {
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
func (repo *ReplyRepoCRUD) Delete(rid uint64) error {
	var (
		res          sql.Result
		rowsAffected int64
		err          error
	)
	if res, err = repo.DB.Exec("DELETE FROM replies WHERE id = ?", rid); err != nil {
		return err
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected > 0 {
		return nil
	}
	return errors.New("failed to delete the reply")
}

func (repo *ReplyRepoCRUD) DeleteGroup(pid uint64) error {
	var (
		res          sql.Result
		rowsAffected int64
		err          error
	)
	if res, err = repo.DB.Exec("DELETE FROM replies WHERE post_id_fkey = ?", pid); err != nil {
		return err
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected > 0 {
		return nil
	}
	return errors.New("failed to delete replies for the post")
}

func (repo *ReplyRepoCRUD) CountReplies(pid uint64) (replies string, err error) {
	if err = repo.DB.QueryRow("SELECT count(id) FROM replies WHERE post_id_fkey = ?", pid).Scan(&replies); err != nil {
		if err != sql.ErrNoRows {
			return "0", err
		}
		return "0", nil
	}
	return replies, nil
}
