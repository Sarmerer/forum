package repository

import "forum/api/models"

type PostRepo interface {
	FindAll() ([]models.Post, error)
	FindByID(uint64) (*models.Post, int, error)
	Create(*models.Post) error
	Update(*models.Post) error
	Delete(uint64) (int, error)
}
