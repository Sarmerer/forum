package controllers

import (
	"net/http"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
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

func GetCategoriesByPostID(pid int64) ([]models.Category, error) {
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

func UpdateCategories(pid int64, categories []string) error {
	var (
		repo repository.CategoryRepo = crud.NewCategoryRepoCRUD()
		err  error
	)
	if err = repo.DeleteGroup(pid); err != nil {
		return err
	}
	if err = repo.Create(pid, categories); err != nil {
		return err
	}
	return nil
}

func DeleteAllCategoriesForPost(pid int64) error {
	var (
		repo repository.CategoryRepo = crud.NewCategoryRepoCRUD()
		err  error
	)
	if err = repo.DeleteGroup(pid); err != nil {
		return err
	}
	return nil
}