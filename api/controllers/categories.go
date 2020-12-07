package controllers

import (
	"net/http"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
)

// GetAllCategories returns categories for all posts as an array
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	var (
		repo       repository.CategoryRepo = crud.NewCategoryRepoCRUD()
		categories []models.Category
		err        error
	)
	if categories, err = repo.FindAll(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}
	response.Success(w, nil, categories)
}
