package models

import (
	"database/sql"
	"errors"
)

//Category struct contains info about category
type Category struct {
	ID          int
	Name        string
	Description string
}

//CategoryModel helps performing CRUD operations
type CategoryModel struct {
	DB *sql.DB
}

//NewCategoryModel creates an instance of CategoryModel
func NewCategoryModel(db *sql.DB) (*CategoryModel, error) {
	return &CategoryModel{db}, nil
}

//FindAll returns all categories in the database
func (um *CategoryModel) FindAll() ([]Category, error) {
	rows, e := um.DB.Query("SELECT * FROM categories")
	if e != nil {
		return nil, e
	}
	var categories []Category

	for rows.Next() {
		var category Category
		rows.Scan(&category.ID, &category.Name, &category.Description)
		categories = append(categories, category)
	}
	return categories, nil
}

//Find returns a specific category from the database
func (um *CategoryModel) Find(id int) (Category, error) {
	var category Category
	rows, err := um.DB.Query("SELECT * FROM categories WHERE category_id = ?", id)
	if err != nil {
		return category, err
	}
	for rows.Next() {
		rows.Scan(&category.ID, &category.Name, &category.Description)
	}
	return category, nil
}

//Create adds a new category to the database
func (um *CategoryModel) Create(category *Category) error {
	statement, err := um.DB.Prepare("INSERT INTO categories (category_name, category_description) VALUES (?, ?)")
	if err != nil {
		return err
	}
	res, err := statement.Exec(category.Name, category.Description)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("could not create new category")
}

//Update updates existing category in the database
func (um *CategoryModel) Update(category *Category) error {
	statement, err := um.DB.Prepare("UPDATE categories SET category_name = ?, category_description = ? WHERE category_id = ?")
	if err != nil {
		return err
	}
	res, err := statement.Exec(category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("could not update the category")
}

//Delete deletes category from the database
func (um *CategoryModel) Delete(id int) error {
	res, err := um.DB.Exec("DELETE FROM categories WHERE category_id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("could not delete the category")
}
