package repository

import "forum/api/models"

type ReplyRepo interface {
	FindAll(uint64) ([]models.PostReply, error)
	FindByID(uint64) (*models.PostReply, int, error)
	Create(*models.PostReply) error
	Update(*models.PostReply) error
	Delete(uint64) error
	DeleteGroup(uint64) error
}
