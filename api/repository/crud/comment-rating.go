package crud

import (
	"errors"
)

func (CommentRepoCRUD) GetRating(commentID int64, currentUser int64) (int, int, error) {
	var (
		rating       int
		yourReaction int
	)
	return rating, yourReaction, nil
}

func (CommentRepoCRUD) Rate(commentID, userID int64, reaction int) error {

	return errors.New("could not set rating")
}

func (CommentRepoCRUD) DeleteAllReactions(cid int64) error {

	return nil
}
