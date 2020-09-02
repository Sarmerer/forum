package crud

import (
	"database/sql"
	"errors"
	"forum/api/models"
	"forum/api/repository"
	"forum/config"
	"net/http"
	"time"
)

//ReplyRepoCRUD helps performing CRUD operations
type ReplyRepoCRUD struct{}

//NewReplyRepoCRUD creates an instance of PostReplyModel
func NewReplyRepoCRUD() *ReplyRepoCRUD {
	return &ReplyRepoCRUD{}
}

//FindAll returns all replies for the specified post
func (ReplyRepoCRUD) FindAll(postID uint64) ([]models.PostReply, error) {
	var (
		rows    *sql.Rows
		replies []models.PostReply
		err     error
	)
	if rows, err = repository.DB.Query(
		"SELECT * FROM replies WHERE post_id_fkey = ?", postID); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	for rows.Next() {
		var r models.PostReply
		rows.Scan(&r.ID, &r.Author, &r.Content, &r.Created, &r.Post)
		replies = append(replies, r)
	}
	return replies, nil
}

//FindByID returns a specific reply from the database
func (ReplyRepoCRUD) FindByID(rid uint64) (*models.PostReply, int, error) {
	var (
		r   models.PostReply
		err error
	)
	if err = repository.DB.QueryRow(
		"SELECT * FROM replies WHERE id = ?", rid,
	).Scan(
		&r.ID, &r.Author, &r.Content, &r.Created, &r.Post,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusBadRequest, errors.New("reply not found")
	}
	return &r, http.StatusOK, nil
}

//Create adds a new reply to the database
func (ReplyRepoCRUD) Create(r *models.PostReply) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"INSERT INTO replies (author_fkey, content, created, post_id_fkey) VALUES (?, ?, ?, ?)",
		r.Author, r.Content, time.Now().Format(config.TimeLayout), r.Post,
	); err != nil {
		return err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return nil
}

//Update updates existing reply in the database
func (ReplyRepoCRUD) Update(r *models.PostReply) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"UPDATE replies SET author_fkey = ?, content = ?, created = ?, post_id_fkey = ? WHERE id = ?",
		r.Author, r.Content, r.Created, r.Post, r.ID,
	); err != nil {
		return err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return nil
}

//Delete deletes reply from the database
func (ReplyRepoCRUD) Delete(rid uint64) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"DELETE FROM replies WHERE id = ?", rid,
	); err != nil {
		return err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return nil
}

func (ReplyRepoCRUD) DeleteGroup(pid uint64) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"DELETE FROM replies WHERE post_id_fkey = ?", pid,
	); err != nil {
		return err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return nil
}

func (ReplyRepoCRUD) CountReplies(pid uint64) (replies string, err error) {
	if err = repository.DB.QueryRow(
		"SELECT count(id) FROM replies WHERE post_id_fkey = ?", pid,
	).Scan(
		&replies,
	); err != nil {
		if err != sql.ErrNoRows {
			return "0", err
		}
		return "0", nil
	}
	return replies, nil
}
