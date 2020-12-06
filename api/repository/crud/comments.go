package crud

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/utils"
)

//CommentRepoCRUD helps performing CRUD operations
type CommentRepoCRUD struct{}

//NewCommentRepoCRUD creates an instance of PostReplyModel
func NewCommentRepoCRUD() CommentRepoCRUD {
	return CommentRepoCRUD{}
}

func fetchAuthor(c *models.Comment) (err error) {
	var status int
	if c.Author, status, err = NewUserRepoCRUD().FindByID(c.AuthorID); err != nil {
		if status == http.StatusInternalServerError {
			return err
		}
		c.Author = DeletedUser
	}
	return nil
}

//FindAll returns all replies for the specified post
func (CommentRepoCRUD) FindByPostID(postID, userID int64) ([]models.Comment, error) {
	var (
		rows     *sql.Rows
		comments []models.Comment
		err      error
	)
	if rows, err = repository.DB.Query(
		`SELECT *,
		(
			SELECT TOTAL(reaction)
			FROM comments_reactions
			WHERE comment_id_fkey = c.id
		) AS rating,
		IFNULL (
			(
				SELECT reaction
				FROM comments_reactions
				WHERE user_id_fkey = $1
					AND comment_id_fkey = c.id
			),
			0
		) AS yor_reaction
		FROM comments c
		WHERE post_id_fkey = $2
		ORDER BY created DESC`,
		userID, postID,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	for rows.Next() {
		var c models.Comment
		rows.Scan(&c.ID, &c.AuthorID, &c.Content, &c.Created, &c.PostID, &c.Edited, &c.Rating, &c.YourReaction)
		if err = fetchAuthor(&c); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (CommentRepoCRUD) FindByAuthor(userID, requestorID int64) ([]models.Comment, int, error) {
	var (
		rows     *sql.Rows
		comments []models.Comment
		err      error
	)
	if rows, err = repository.DB.Query(
		`SELECT *,
		(
			SELECT TOTAL(reaction)
			FROM comments_reactions
			WHERE comment_id_fkey = c.id
		) AS rating,
		IFNULL (
			(
				SELECT reaction
				FROM comments_reactions
				WHERE user_id_fkey = $1
					AND comment_id_fkey = c.id
			),
			0
		) AS yor_reaction
		FROM comments c
		WHERE author_id_fkey = $2
		ORDER BY created DESC`,
		requestorID, userID,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, nil
	}
	for rows.Next() {
		var c models.Comment
		rows.Scan(&c.ID, &c.AuthorID, &c.Content, &c.Created, &c.PostID, &c.Edited, &c.Rating, &c.YourReaction)
		if err = fetchAuthor(&c); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		comments = append(comments, c)
	}
	return comments, http.StatusOK, nil
}

//FindByID returns a specific reply from the database
func (CommentRepoCRUD) FindByID(commentID int64) (*models.Comment, int, error) {
	var (
		c   models.Comment
		err error
	)
	if err = repository.DB.QueryRow(
		`SELECT *
		FROM comments
		WHERE id = ?`, commentID,
	).Scan(
		&c.ID, &c.AuthorID, &c.Content, &c.Created, &c.PostID, &c.Edited,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusBadRequest, errors.New("reply not found")
	}
	if err = fetchAuthor(&c); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &c, http.StatusOK, nil
}

//Create adds a new reply to the database
func (CommentRepoCRUD) Create(comment *models.Comment) (*models.Comment, error) {
	var (
		result       sql.Result
		now          int64 = utils.CurrentTime()
		newComment   *models.Comment
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`INSERT INTO comments (
			author_id_fkey,
			content,
			created,
			post_id_fkey,
			edited
		)
		VALUES (?, ?, ?, ?, ?)`,
		comment.AuthorID, comment.Content, now, comment.PostID, now,
	); err != nil {
		return nil, err
	}
	if comment.ID, err = result.LastInsertId(); err != nil {
		return nil, err
	}
	if newComment, _, err = NewCommentRepoCRUD().FindByID(comment.ID); err != nil {
		return nil, err
	}
	if err = fetchAuthor(newComment); err != nil {
		return nil, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return nil, err
	}
	if rowsAffected > 0 {
		return newComment, err
	}
	return nil, errors.New("could not create a comment")
}

//Update updates existing reply in the database
func (CommentRepoCRUD) Update(comment *models.Comment) (*models.Comment, error) {
	var (
		result         sql.Result
		updatedComment *models.Comment
		rowsAffected   int64
		err            error
	)
	if result, err = repository.DB.Exec(
		`UPDATE comments
		SET author_id_fkey = ?,
			content = ?,
			created = ?,
			post_id_fkey = ?,
			edited = ?
		WHERE id = ?`,
		comment.AuthorID, comment.Content, comment.Created, comment.PostID, utils.CurrentTime(), comment.ID,
	); err != nil {
		return nil, err
	}
	if updatedComment, _, err = NewCommentRepoCRUD().FindByID(comment.ID); err != nil {
		return nil, err
	}
	if err = fetchAuthor(updatedComment); err != nil {
		return nil, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return nil, err
	}
	if rowsAffected > 0 {
		return updatedComment, nil
	}
	return nil, errors.New("couldn't update the comment")
}

//Delete deletes reply from the database
func (CommentRepoCRUD) Delete(cid int64) error {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
	)
	ctx = context.Background()
	tx, err = repository.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.ExecContext(ctx,
		`DELETE FROM comments_reactions
		WHERE comment_id_fkey = $1`, cid)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx,
		`DELETE FROM comments
		WHERE id = $1`, cid)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (CommentRepoCRUD) DeleteGroup(pid int64) error {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
	)
	ctx = context.Background()
	tx, err = repository.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		`DELETE FROM comments_reactions
			WHERE comment_id_fkey IN (
				SELECT id
				FROM comments
				WHERE post_id_fkey = $1
			)`, pid)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx,
		`DELETE FROM comments
		WHERE post_id_fkey = $1`, pid)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (CommentRepoCRUD) Count(pid int64) (comments string, err error) {
	if err = repository.DB.QueryRow(
		`SELECT count(id)
		FROM comments
		WHERE post_id_fkey = ?`, pid,
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
