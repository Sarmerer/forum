package crud

import (
	"errors"
	"net/http"

	"github.com/sarmerer/forum/api/models"
)

//CommentRepoCRUD helps performing CRUD operations
type CommentRepoCRUD struct{}

//NewCommentRepoCRUD creates an instance of PostReplyModel
func NewCommentRepoCRUD() CommentRepoCRUD {
	return CommentRepoCRUD{}
}

//FindAll returns all replies for the specified post
func (CommentRepoCRUD) FindByPostID(postID, userID int64) ([]models.Comment, error) {
	var (
		comments []models.Comment
		err      error
	)
	return comments, err
}

func (CommentRepoCRUD) FindByAuthor(userID, requestorID int64) ([]models.Comment, int, error) {
	var (
		comments []models.Comment
		err      error
	)
	return comments, http.StatusOK, err
}

//FindByID returns a specific reply from the database
func (CommentRepoCRUD) FindByID(commentID int64) (*models.Comment, int, error) {
	var (
		c   models.Comment
		err error
	)
	return &c, http.StatusOK, err
}

//Create adds a new reply to the database
func (CommentRepoCRUD) Create(comment *models.Comment) (*models.Comment, error) {
	return nil, errors.New("could not create a comment")
}

//Update updates existing reply in the database
func (CommentRepoCRUD) Update(comment *models.Comment) (*models.Comment, error) {

	return nil, errors.New("couldn't update the comment")
}

//Delete deletes reply from the database
func (CommentRepoCRUD) Delete(cid int64) error {

	return nil
}

func (CommentRepoCRUD) DeleteGroup(pid int64) error {

	return nil
}

func (CommentRepoCRUD) Count(pid int64) (comments string, err error) {

	return comments, nil
}
