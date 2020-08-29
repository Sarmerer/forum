package repository

import "forum/api/models"

type ReplyRepo interface {
	FindAll(uint64) ([]models.PostReply, error)
	//TODO: add int to return values
	FindByID(int) (models.PostReply, error)
	Create(*models.PostReply) (int, error)
	//TODO: change bool to error
	Update(*models.PostReply) bool
	Delete(int) bool
}
