package repository

import (
	"forum/api/models"
)

type UserRepo interface {
	FindAll() ([]models.User, error)
	FindByID(uint64) (*models.User, int, error)
	Create(*models.User) (int, error)
	Update(*models.User) (int, error)
	Delete(uint64) (int, error)
	FindByNameOrEmail(string) (*models.User, int, error)

	UpdateSession(uint64, string) error
	ValidateSession(string) (uint64, error)

	GetRole(uint64) (int, int, error)
	UpdateRole(uint64, int) error
}
