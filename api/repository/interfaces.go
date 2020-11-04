package repository

import (
	"github.com/sarmerer/forum/api/models"
)

type UserRepo interface {
	FindAll() ([]models.User, error)
	FindByID(int64) (*models.User, int, error)
	Create(*models.User) (int, error)
	Update(*models.User) (int, error)
	Delete(int64) (int, error)
	FindByNameOrEmail(string) (*models.User, int, error)

	UpdateSession(int64, string) error
	ValidateSession(string) (models.UserCtx, int, error)

	GetRole(int64) (int, int, error)
	UpdateRole(int64, int) error
}

type PostRepo interface {
	// TODO retourn your_reaction in all of the following methods
	// FindAll takes in user id, current page and offset,
	// user id for server to know if that user put a like or dislike on a post,
	// current page and offset for pagination
	FindAll(userID int64, perPage, currentPage int) (*models.Posts, error)
	// FindByID takes in post id and user id, for the same reasaon, as in FindAll()
	FindByID(postID, userID int64) (*models.Post, int, error)
	FindByAuthor(userID int64) ([]models.Post, error)
	FindByCategories(categories []string) ([]models.Post, error)

	Create(*models.Post) (int64, error)
	Update(*models.Post) error
	Delete(postID int64) (int, error)

	Rate(postID, userID int64, reaction int) error
	GetRating(postID, userID int64) (int, int, error)
	// TODO delete all reactions from user, when deleting a user
	DeleteAllReactions(postID int64) error
}

type CategoryRepo interface {
	FindAll() ([]models.Category, error)
	FindByPostID(postID int64) ([]models.Category, error)
	Create(postID int64, categories []string) error
	DeleteGroup(postID int64) error
}

type CommentRepo interface {
	FindAll(postID, userID int64) ([]models.Comment, error)
	FindByID(int64) (*models.Comment, int, error)
	Create(*models.Comment) error
	Update(*models.Comment) error
	Delete(int64) error
	DeleteGroup(int64) error
	Count(int64) (string, error)

	Rate(int64, int64, int) error
	GetRating(int64, int64) (int, int, error)
	DeleteAllReactions(int64) error
}
