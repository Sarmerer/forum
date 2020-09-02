package controllers

import (
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
)

func GetCategories(pid uint64) ([]models.Category, error) {
	var (
		repo       repository.CategoryRepo = crud.NewCategoryRepoCRUD()
		categories []models.Category
		err        error
	)
	if categories, err = repo.FindAll(pid); err != nil {
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
