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
