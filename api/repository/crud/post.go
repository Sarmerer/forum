package crud

import (
	"errors"
	"net/http"

	"github.com/sarmerer/forum/api/models"
)

//PostRepoCRUD helps performing CRUD operations
type PostRepoCRUD struct{}

//NewPostRepoCRUD creates an instance of PostModel
func NewPostRepoCRUD() PostRepoCRUD {
	return PostRepoCRUD{}
}

//FindAll returns all posts in the database
//TODO improve this MESS
func (PostRepoCRUD) FindAll(userID int64, input models.InputAllPosts) (*models.Posts, int, error) {
	var (
		result models.Posts
		err    error
	)
	return &result, http.StatusOK, err
}

func (PostRepoCRUD) FindRecent(amount int) ([]models.Post, int, error) {
	var (
		result []models.Post
		err    error
	)

	return result, http.StatusOK, err
}

//FindByID returns a specific post from the database
func (PostRepoCRUD) FindByID(postID int64, userID int64) (*models.Post, int, error) {
	var (
		p   models.Post
		err error
	)

	return &p, http.StatusOK, err
}

func (PostRepoCRUD) FindByAuthor(userID, requestorID int64) ([]models.Post, int, error) {
	var (
		posts []models.Post
		err   error
	)

	return posts, http.StatusOK, err
}

func (PostRepoCRUD) FindByCategories(categories []string, requestorID int64) ([]models.Post, int, error) {
	var (
		posts []models.Post
		err   error
	)
	return posts, http.StatusOK, err
}

//Create adds a new post to the database
func (PostRepoCRUD) Create(post *models.Post, categories []string) (*models.Post, int, error) {
	return nil, http.StatusBadRequest, errors.New("could not create the post")
}

//Update updates existing post in the database
func (PostRepoCRUD) Update(post *models.Post, userCtx models.UserCtx) (*models.Post, int, error) {

	return nil, http.StatusBadRequest, errors.New("could not update the post")
}

//Delete deletes post from the database
func (PostRepoCRUD) Delete(pid int64) (int, error) {

	return http.StatusNotModified, errors.New("could not delete the post")
}
