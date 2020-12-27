package crud

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
)

func (CommentRepoCRUD) softDelete(commentID int64) error {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
	)
	ctx = context.Background()
	if tx, err = repository.DB.BeginTx(ctx, nil); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx,
		`UPDATE comments
		SET deleted = ?,
			author_id_fkey = ?,
			content = ?,
			edited = ?
		WHERE _id = ?`,
		true, 0, "", 0, commentID,
	); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx,
		`DELETE FROM comments_reactions
		WHERE comment_id_fkey = ?`, commentID,
	); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (CommentRepoCRUD) hardDelete(comment *models.Comment) error {
	var (
		hasChildren   bool
		performDelete func(*models.Comment) error
		err           error
	)
	performDelete = func(comment *models.Comment) error {
		var (
			ctx context.Context
			tx  *sql.Tx
		)
		ctx = context.Background()
		if tx, err = repository.DB.BeginTx(ctx, nil); err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx,
			`DELETE FROM comments_reactions
			WHERE comment_id_fkey = ?`, comment.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.ExecContext(ctx,
			`DELETE FROM comments
			WHERE _id = ?`, comment.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
		if comment.ParentID != 0 {
			if hasChildren, err = NewCommentRepoCRUD().hasChildren(comment.ParentID); err != nil {
				return err
			}
			if !hasChildren {
				var parent *models.Comment
				if parent, _, err = NewCommentRepoCRUD().FindByID(comment.ParentID); err != nil {
					return err
				}
				fmt.Println("deleting", parent.ID)
				performDelete(parent)
			}
		}
		return nil
	}
	if err = performDelete(comment); err != nil {
		return err
	}
	return nil
}

func (CommentRepoCRUD) hasChildren(commentID int64) (bool, error) {
	var (
		temp int64
		err  error
	)
	if err = repository.DB.QueryRow("SELECT COUNT(_id) FROM comments WHERE parent_id_fkey = ? AND deleted = 0",
		commentID,
	).Scan(&temp); err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	if temp > 0 {
		return true, nil
	}
	return false, nil
}
