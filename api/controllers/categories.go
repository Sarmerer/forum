package controllers

import (
	"forum/api/helpers"
	"forum/api/models"
	"forum/api/repository"
)

func GetCategories(pid uint64) ([]models.Category, error) {
	var (
		cm         repository.CategoryRepo
		categories []models.Category
		err        error
	)

	if cm, err = helpers.PrepareCategoriesRepo(); err != nil {
		return nil, err
	}
	if categories, err = cm.FindAll(pid); err != nil {
		return nil, err
	}
	return categories, nil
}

func DeleteAllCategoriesForPost(pid uint64) error {
	var (
		cm  repository.CategoryRepo
		err error
	)
	if cm, err = helpers.PrepareCategoriesRepo(); err != nil {
		return err
	}
	if err = cm.DeleteGroup(pid); err != nil {
		return err
	}
	return nil
}
