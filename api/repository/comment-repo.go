package repository

import "forum/api/models"

type CommentRepo interface {
	FindAll(int64) ([]models.PostComment, error)
	FindByID(int64) (*models.PostComment, int, error)
	Create(*models.PostComment) error
	Update(*models.PostComment) error
	Delete(int64) error
	DeleteGroup(int64) error
	Count(int64) (string, error)
}
