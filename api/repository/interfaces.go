package repository

import (
	"forum/api/models"
)

type UserRepo interface {
	FindAll() ([]models.User, error)
	FindByID(int64) (*models.User, int, error)
	Create(*models.User) (int, error)
	Update(*models.User) (int, error)
	Delete(int64) (int, error)
	FindByNameOrEmail(string) (*models.User, int, error)

	UpdateSession(int64, string) error
	ValidateSession(string) (int64, int, error)

	GetRole(int64) (int, int, error)
	UpdateRole(int64, int) error
}

type PostRepo interface {
	FindAll(int64) ([]models.Post, error)
	FindByID(int64, int64) (*models.Post, int, error)
	FindByAuthor(int64) ([]models.Post, error)
	FindByCategories([]string) ([]models.Post, error)

	Create(*models.Post) (int64, error)
	Update(*models.Post) error
	Delete(int64) (int, error)

	RatePost(int64, int64, int) error
	GetRating(int64, int64) (int, int, error)
	DeleteAllReactions(int64) error
}

type CategoryRepo interface {
	FindAll() ([]models.Category, error)
	FindByPostID(int64) ([]models.Category, error)
	Create(int64, []string) error
	DeleteGroup(int64) error
}

type CommentRepo interface {
	FindAll(int64) ([]models.Comment, error)
	FindByID(int64) (*models.Comment, int, error)
	Create(*models.Comment) error
	Update(*models.Comment) error
	Delete(int64) error
	DeleteGroup(int64) error
	Count(int64) (string, error)
}
