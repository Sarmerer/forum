package crud

import (
	"errors"
)

func (PostRepoCRUD) GetRating(postID int64, currentUser int64) (int, int, error) {
	var (
		rating       int
		yourReaction int
		err          error
	)
	return rating, yourReaction, err
}

func (PostRepoCRUD) Rate(postID, userID int64, reaction int) error {
	return errors.New("could not set rating")
}

func (PostRepoCRUD) DeleteAllReactions(pid int64) error {
	return nil
}
