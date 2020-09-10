package controllers

import (
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"net/http"
)

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

func GetCategoriesByPostID(pid uint64) ([]models.Category, error) {
	var (
		repo       repository.CategoryRepo = crud.NewCategoryRepoCRUD()
		categories []models.Category
		err        error
	)
	if categories, err = repo.FindByPostID(pid); err != nil {
		return nil, err
	}
	return categories, nil
}

func DeleteAllCategoriesForPost(pid uint64) error {
	var (
		repo repository.CategoryRepo = crud.NewCategoryRepoCRUD()
		err  error
	)
	if err = repo.DeleteGroup(pid); err != nil {
		return err
	}
	return nil
}
