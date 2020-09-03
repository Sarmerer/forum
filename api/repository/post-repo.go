package repository

import "forum/api/models"

type PostRepo interface {
	FindAll() ([]models.Post, error)
	FindByID(uint64) (*models.Post, int, error)
	FindByCategories(string) ([]models.Post, error)
	Create(*models.Post) (int64, error)
	Update(*models.Post) error
	Delete(uint64) (int, error)
}
