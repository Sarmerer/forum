package repository

import "forum/api/models"

type CategoryRepo interface {
	FindAll() ([]models.Category, error)
	FindByPostID(uint64) ([]models.Category, error)
	Create(int64, ...string) error
	DeleteGroup(uint64) error
}
