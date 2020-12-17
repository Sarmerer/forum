package crud

import (
	"database/sql"
	"errors"

	"github.com/sarmerer/forum/api/repository"
)

func (CommentRepoCRUD) GetRating(commentID int64, currentUser int64) (int, int, error) {
	var (
		rating       int
		yourReaction int
		err          error
	)
	if err = repository.DB.QueryRow(
		`SELECT TOTAL(reaction) AS rating,
		IFNULL (
			(
				SELECT reaction
				FROM comments_reactions
				WHERE user_id_fkey = $1
					AND comment_id_fkey = $2
			),
			0
		) AS yor_reaction
		FROM comments_reactions
		WHERE comment_id_fkey = $2`,
		currentUser, commentID,
	).Scan(&rating, &yourReaction); err != nil {
		if err != sql.ErrNoRows {
			return 0, 0, err
		}
		return 0, 0, errors.New("comment not found")
	}
	return rating, yourReaction, nil
}

func (CommentRepoCRUD) Rate(commentID, userID int64, reaction int) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	switch reaction {
	case 0:
		if result, err = repository.DB.Exec(
			`DELETE FROM comments_reactions
			WHERE comment_id_fkey = ?
				AND user_id_fkey = ?`,
			commentID, userID,
		); err != nil {
			return err
		}
	default:
		if result, err = repository.DB.Exec(
			`INSERT
			OR REPLACE INTO comments_reactions(_id, comment_id_fkey, user_id_fkey, reaction)
		VALUES (
				(
					SELECT _id
					FROM comments_reactions
					WHERE comment_id_fkey = $1
						AND user_id_fkey = $2
				),
				$1,
				$2,
				$3
			);`,
			commentID, userID, reaction,
		); err != nil {
			return err
		}
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("could not set rating")
}

func (CommentRepoCRUD) DeleteAllReactions(cid int64) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)

	if result, err = repository.DB.Exec(
		`DELETE FROM comments_reactions
		WHERE comment_id_fkey = ?`, cid,
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
