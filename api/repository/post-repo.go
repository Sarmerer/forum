package repository

import "forum/api/models"

type PostRepo interface {
	FindAll() ([]models.Post, error)
	FindByID(int64) (*models.Post, int, error)
	FindByAuthor(int64) ([]models.Post, error)
	FindByCategories([]string) ([]models.Post, error)

	Create(*models.Post) (int64, error)
	Update(*models.Post) error
	Delete(int64) (int, error)

	RatePost(int64, int64, int) error
	GetRating(int64, int64) (int, int, error)
}
