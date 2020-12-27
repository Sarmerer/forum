package crud

import (
	"context"
	"database/sql"
	"errors"
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
func (CommentRepoCRUD) FindByPostID(postID, userID int64) ([]*models.Comment, error) {
	var (
		rows     *sql.Rows
		comments []*models.Comment
		cache    map[int64]*models.Comment = make(map[int64]*models.Comment)

		increment func(map[int64]*models.Comment, int64)
		err       error
	)
	if rows, err = repository.DB.Query(
		`SELECT *,
		(
			SELECT TOTAL(reaction)
			FROM comments_reactions
			WHERE comment_id_fkey = c._id
		) AS rating,
		IFNULL (
			(
				SELECT reaction
				FROM comments_reactions
				WHERE user_id_fkey = $1
					AND comment_id_fkey = c._id
			),
			0
		) AS your_reaction
		FROM comments c
		WHERE post_id_fkey = $2
		ORDER BY lineage, rating DESC`,
		userID, postID,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	increment = func(c map[int64]*models.Comment, parentID int64) {
		c[parentID].ChildrenLen++
		if c[parentID].ParentID != 0 {
			increment(c, c[parentID].ParentID)
		}
	}
	for rows.Next() {
		var c models.Comment
		rows.Scan(&c.ID, &c.AuthorID, &c.PostID, &c.ParentID, &c.Depth, &c.Lineage,
			&c.Content, &c.Created, &c.Edited, &c.Deleted, &c.Rating, &c.YourReaction)
		if err = fetchAuthor(&c); err != nil {
			return nil, err
		}
		if c.ParentID != 0 {
			// fmt.Printf("id: %d, parent: %d, cache: %v\n", c.ID, c.ParentID, cache[c.ParentID] != nil)
			cache[c.ParentID].Children = append(cache[c.ParentID].Children, &c)
			cache[c.ID] = &c
			increment(cache, c.ParentID)
		} else {
			comments = append(comments, &c)
			cache[c.ID] = comments[len(comments)-1]
			cache[c.ID].ChildrenLen++
		}
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
			WHERE comment_id_fkey = c._id
		) AS rating,
		IFNULL (
			(
				SELECT reaction
				FROM comments_reactions
				WHERE user_id_fkey = $1
					AND comment_id_fkey = c._id
			),
			0
		) AS yor_reaction
		FROM comments c
		WHERE author_id_fkey = $2
		AND deleted = 0
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
		rows.Scan(&c.ID, &c.AuthorID, &c.PostID, &c.ParentID, &c.Depth, &c.Lineage,
			&c.Content, &c.Created, &c.Edited, &c.Deleted, &c.Rating, &c.YourReaction)
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
		WHERE _id = ?`, commentID,
	).Scan(
		&c.ID, &c.AuthorID, &c.PostID, &c.ParentID, &c.Depth,
		&c.Lineage, &c.Content, &c.Created, &c.Edited, &c.Deleted,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusBadRequest, errors.New("comment not found")
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
		now          int64 = utils.CurrentUnixTime()
		newComment   *models.Comment
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`INSERT INTO comments (
			author_id_fkey,
			post_id_fkey,
			parent_id_fkey,
			depth,
			lineage,
			content,
			created,
			edited,
			deleted
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		comment.AuthorID, comment.PostID, comment.ParentID,
		comment.Depth, comment.Lineage, comment.Content, now, now, false,
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
			post_id_fkey = ?,
			parent_id_fkey = ?,
			depth = ?,
			lineage = ?,
			content = ?,
			created = ?,
			edited = ?,
			deleted = ?
		WHERE _id = ?`,
		comment.AuthorID, comment.PostID, comment.ParentID, comment.Depth, comment.Lineage,
		comment.Content, comment.Created, utils.CurrentUnixTime(), false, comment.ID,
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

// Delete deletes reply from the database
func (CommentRepoCRUD) Delete(comment *models.Comment) error {
	var (
		hasChildren bool
		err         error
	)

	if hasChildren, err = NewCommentRepoCRUD().hasChildren(comment.ID); err != nil {
		return err
	}
	if hasChildren {
		if err = NewCommentRepoCRUD().softDelete(comment.ID); err != nil {
			return err
		}
		return nil
	}
	if err = NewCommentRepoCRUD().hardDelete( comment); err != nil {
		return err
	}
	return nil
}

func (CommentRepoCRUD) DeleteGroup(postID int64) error {
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
				SELECT _id
				FROM comments
				WHERE post_id_fkey = ?
			)`, postID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx,
		`DELETE FROM comments
		WHERE post_id_fkey = $1`, postID)
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

func (CommentRepoCRUD) CountForPost(post *models.Post) error {
	var err error
	if err = repository.DB.QueryRow(
		`SELECT COUNT(_id) AS total_comments,
		COUNT(DISTINCT author_id_fkey) AS total_participants
		FROM comments
		WHERE post_id_fkey = $1
		AND deleted = 0`, post.ID,
	).Scan(
		&post.CommentsCount, &post.ParticipantsCount,
	); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return nil
	}
	return nil
}
