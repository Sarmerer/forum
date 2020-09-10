package repository

import "forum/api/models"

type PostRepo interface {
	FindAll() ([]models.Post, error)
	FindByID(uint64) (*models.Post, int, error)
	FindByAuthor(uint64) ([]models.Post, error)
	FindByCategories([]string) ([]models.Post, error)

	Create(*models.Post) (int64, error)
	Update(*models.Post) error
	Delete(uint64) (int, error)

	RatePost(uint64, uint64, int) error
	GetRating(uint64) (int, error)
}
