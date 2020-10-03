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

//CommentRepoCRUD helps performing CRUD operations
type CommentRepoCRUD struct{}

//NewCommentRepoCRUD creates an instance of PostReplyModel
func NewCommentRepoCRUD() *CommentRepoCRUD {
	return &CommentRepoCRUD{}
}

//FindAll returns all replies for the specified post
func (CommentRepoCRUD) FindAll(postID int64) ([]models.PostComment, error) {
	var (
		rows    *sql.Rows
		replies []models.PostComment
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
		var r models.PostComment
		rows.Scan(&r.ID, &r.AuthorID, &r.AuthorName, &r.Content, &r.Created, &r.Post, &r.Edited)
		replies = append(replies, r)
	}
	return replies, nil
}

//FindByID returns a specific reply from the database
func (CommentRepoCRUD) FindByID(rid int64) (*models.PostComment, int, error) {
	var (
		r   models.PostComment
		err error
	)
	if err = repository.DB.QueryRow(
		"SELECT * FROM replies WHERE id = ?", rid,
	).Scan(
		&r.ID, &r.AuthorID, &r.AuthorName, &r.Content, &r.Created, &r.Post, &r.Edited,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusBadRequest, errors.New("reply not found")
	}
	return &r, http.StatusOK, nil
}

//Create adds a new reply to the database
func (CommentRepoCRUD) Create(r *models.PostComment) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"INSERT INTO replies (author_fkey, author_name_fkey, content, created, post_id_fkey, edited) VALUES (?, ?, ?, ?, ?, ?)",
		r.AuthorID, r.AuthorName, r.Content, time.Now().Format(config.TimeLayout), r.Post, 0,
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
func (CommentRepoCRUD) Update(r *models.PostComment) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"UPDATE replies SET author_fkey = ?, author_name_fkey = ?, content = ?, created = ?, post_id_fkey = ?, edited = ? WHERE id = ?",
		r.AuthorID, r.AuthorName, r.Content, r.Created, r.Post, 1, r.ID,
	); err != nil {
		return err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("couldn't update the comment")
}

//Delete deletes reply from the database
func (CommentRepoCRUD) Delete(rid int64) error {
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

func (CommentRepoCRUD) DeleteGroup(pid int64) error {
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

func (CommentRepoCRUD) Count(pid int64) (replies string, err error) {
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
