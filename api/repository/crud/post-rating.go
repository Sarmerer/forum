package crud

import (
	"database/sql"
	"errors"
	"forum/api/repository"
)

func (PostRepoCRUD) GetRating(postID int64, currentUser int64) (int, int, error) {
	var (
		rating       int
		yourReaction int
		err          error
	)
	if err = repository.DB.QueryRow(
		`SELECT TOTAL(reaction) AS rating, IFNULL ((SELECT reaction FROM reactions WHERE user_id_fkey = $1 AND post_id_fkey = $2), 0) AS yor_reaction FROM reactions WHERE post_id_fkey = $2`,
		currentUser, postID,
	).Scan(&rating, &yourReaction); err != nil {
		if err != sql.ErrNoRows {
			return 0, 0, err
		}
		return 0, 0, errors.New("post not found")
	}
	return rating, yourReaction, nil
}

func (PostRepoCRUD) RatePost(postID, userID int64, reaction int) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	switch reaction {
	case 0:
		if result, err = repository.DB.Exec(
			`DELETE FROM reactions WHERE post_id_fkey = ? AND user_id_fkey = ?`,
			postID, userID,
		); err != nil {
			return err
		}
	default:
		if result, err = repository.DB.Exec(
			`INSERT OR REPLACE
			INTO reactions(id, post_id_fkey, user_id_fkey, reaction)
			VALUES ((SELECT id FROM reactions WHERE post_id_fkey = $1 AND user_id_fkey = $2), $1, $2, $3);`,
			postID, userID, reaction,
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

// func (PostRepoCRUD) UpdateRating(upORdown string, pid uint64) error {
// 	var (
// 		result sql.Result
// 		err    error
// 	)
// 	switch upORdown {
// 	case "up":
// 		if result, err = repository.DB.Exec(
// 			"UPDATE posts SET rating = rating + 1 WHERE id = ?",
// 			pid,
// 		); err != nil {

// 		}
// 	case "down":
// 		if result, err = repository.DB.Exec(
// 			"UPDATE posts SET rating = rating - 1 WHERE id = ?",
// 			pid,
// 		); err != nil {

// 		}
// 	default:
// 	}
// 	return nil
// }
