package crud

import (
	"net/http"

	"github.com/sarmerer/forum/api/models"
)

//CategoryRepoCRUD helps performing CRUD operations
type CategoryRepoCRUD struct{}

//NewCategoryRepoCRUD creates an instance of CategoryModel
func NewCategoryRepoCRUD() CategoryRepoCRUD {
	return CategoryRepoCRUD{}
}

func (CategoryRepoCRUD) FindAll() ([]models.Category, error) {
	var (
		categories []models.Category
	)

	return categories, nil
}

//FindByPostID returns all categories belonging to  a post
func (CategoryRepoCRUD) FindByPostID(postID int64) ([]models.Category, int, error) {
	var (
		categories []models.Category
	)
	return categories, http.StatusOK, nil
}

//Find returns a specific category from the database
// ? do we even need that
func (CategoryRepoCRUD) Find(id int) (*models.Category, error) {
	var category models.Category
	return &category, nil
}

//Create adds a new category to the database
//FIXME category being duplicated, when creating new post with existing category
func (CategoryRepoCRUD) Create(postID int64, categories []string) (err error) {

	return nil
}

func (CategoryRepoCRUD) Update(postID int64, categories []string) error {

	return nil
}

//Delete deletes category from the database
func (CategoryRepoCRUD) Delete(categotyID int) error {

	return nil
}

func (CategoryRepoCRUD) DeleteGroup(pid int64) error {

	return nil
}
