package crud

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
)

//CommentRepoCRUD helps performing CRUD operations
type CommentRepoCRUD struct{}

//NewCommentRepoCRUD creates an instance of PostReplyModel
func NewCommentRepoCRUD() *CommentRepoCRUD {
	return &CommentRepoCRUD{}
}

//FindAll returns all replies for the specified post
func (CommentRepoCRUD) FindAll(postID int64) ([]models.Comment, error) {
	var (
		rows     *sql.Rows
		comments []models.Comment
		err      error
	)
	if rows, err = repository.DB.Query(
		"SELECT * FROM comments WHERE post_id_fkey = ? ORDER BY created DESC", postID); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	for rows.Next() {
		var r models.Comment
		rows.Scan(&r.ID, &r.AuthorID, &r.AuthorName, &r.Content, &r.Created, &r.PostID, &r.Edited)
		comments = append(comments, r)
	}
	return comments, nil
}

//FindByID returns a specific reply from the database
func (CommentRepoCRUD) FindByID(rid int64) (*models.Comment, int, error) {
	var (
		r      models.Comment
		edited int
		err    error
	)
	if err = repository.DB.QueryRow(
		"SELECT * FROM comments WHERE id = ?", rid,
	).Scan(
		&r.ID, &r.AuthorID, &r.AuthorName, &r.Content, &r.Created, &r.PostID, &edited,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusBadRequest, errors.New("reply not found")
	}
	if edited == 1 {
		r.Edited = true
	} else {
		r.Edited = false
	}
	return &r, http.StatusOK, nil
}

//Create adds a new reply to the database
func (CommentRepoCRUD) Create(r *models.Comment) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"INSERT INTO comments (author_id_fkey, author_name_fkey, content, created, post_id_fkey, edited) VALUES (?, ?, ?, ?, ?, ?)",
		r.AuthorID, r.AuthorName, r.Content, time.Now().Format(config.TimeLayout), r.PostID, 0,
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
func (CommentRepoCRUD) Update(r *models.Comment) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"UPDATE comments SET author_id_fkey = ?, author_name_fkey = ?, content = ?, created = ?, post_id_fkey = ?, edited = ? WHERE id = ?",
		r.AuthorID, r.AuthorName, r.Content, r.Created, r.PostID, 1, r.ID,
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
		"DELETE FROM comments WHERE id = ?", rid,
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
		"DELETE FROM comments WHERE post_id_fkey = ?", pid,
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

func (CommentRepoCRUD) Count(pid int64) (comments string, err error) {
	if err = repository.DB.QueryRow(
		"SELECT count(id) FROM comments WHERE post_id_fkey = ?", pid,
	).Scan(
		&comments,
	); err != nil {
		if err != sql.ErrNoRows {
			return "0", err
		}
		return "0", nil
	}
	return comments, nil
}
