package repository

import (
	"github.com/sarmerer/forum/api/models"
)

type UserRepo interface {
	FindAll() (users []models.User, err error)
	FindByID(userID int64) (user *models.User, status int, err error)
	Create(user *models.User) (newUser *models.User, status int, err error)
	Update(user *models.User) (status int, err error)
	Delete(userID int64) (status int, err error)
	FindByLoginOrEmail(login string) (user *models.User, status int, err error)

	GetPassword(userID int64) (password string, status int, err error)
	UpdateSession(userID int64, newSessionID string) error
	ValidateSession(sessionID string) (userCtx models.UserCtx, status int, err error)

	GetRole(userID int64) (role int, status int, err error)
	UpdateRole(userID int64, newRole int) error

	UpdateLastActivity(userID int64) error
}

type PostRepo interface {
	// TODO retourn your_reaction in all of the following methods
	// FindAll takes in user id, current page and offset,
	// user id for server to know if that user put a like or dislike on a post,
	// current page and offset for pagination
	FindAll(requestorID int64, input models.InputAllPosts) (posts *models.Posts, status int, err error)
	// FindByID takes in post id and user id, for the same reasaon, as in FindAll()
	FindByID(postID, requestorID int64) (post *models.Post, status int, err error)
	FindByAuthor(userID, requestorID int64) (posts []models.Post, status int, err error)
	FindByCategories(categories []string, requestorID int64) (posts []models.Post, status int, err error)

	Create(post *models.Post, categories []string) (newPost *models.Post, status int, err error)
	Update(*models.Post, models.UserCtx) (updatedPost *models.Post, status int, err error)
	Delete(postID int64) (int, error)

	Rate(postID, userID int64, reaction int) error
	GetRating(postID, requestorID int64) (int, int, error)
	// TODO delete all reactions from user, when deleting a user
	DeleteAllReactions(postID int64) error
}

type CategoryRepo interface {
	FindAll() ([]models.Category, error)
	FindByPostID(postID int64) (categories []models.Category, status int, err error)
	Create(postID int64, categories []string) error
	Update(postID int64, categories []string) error
	DeleteGroup(postID int64) error
}

type CommentRepo interface {
	FindByID(commentID int64) (*models.Comment, int, error)
	FindByPostID(postID, requestorID int64) ([]models.Comment, error)
	FindByAuthor(userID, requestorID int64) (comments []models.Comment, status int, err error)
	Create(comment *models.Comment) (*models.Comment, error)
	Update(comment *models.Comment) (*models.Comment, error)
	Delete(commentID int64) error
	DeleteGroup(postID int64) error
	Count(postID int64) (string, error)

	Rate(commentID, userID int64, reaction int) error
	GetRating(commentID int64, userID int64) (rating, userReaction int, err error)
	DeleteAllReactions(postID int64) error
}
