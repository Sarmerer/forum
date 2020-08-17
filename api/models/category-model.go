package models

import (
	"database/sql"
	"forum/api/entities"
)

//CategoryModel helps performing CRUD operations
type CategoryModel struct {
	DB *sql.DB
}

//NewCategoryModel creates an instance of CategoryModel
func NewCategoryModel(db *sql.DB) (*CategoryModel, error) {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "categories" (
	"category_id"	INTEGER,
	"category_name"	TEXT UNIQUE,
	"category_description"	TEXT,
	PRIMARY KEY("category_id")
)`)
	if err != nil {
		return nil, err
	}
	statement.Exec()
	return &CategoryModel{
		DB: db,
	}, nil
}

//FindAll returns all categories in the database
func (um *CategoryModel) FindAll() ([]entities.Category, error) {
	rows, e := um.DB.Query("SELECT * FROM categories")
	if e != nil {
		return nil, e
	}
	var categories []entities.Category

	for rows.Next() {
		var category entities.Category
		rows.Scan(&category.ID, &category.Name, &category.Description)
		categories = append(categories, category)
	}
	return categories, nil
}

//Find returns a specific category from the database
func (um *CategoryModel) Find(id int) (entities.Category, error) {
	var category entities.Category
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
func (um *CategoryModel) Create(category *entities.Category) (bool, string) {
	statement, err := um.DB.Prepare("INSERT INTO categories (category_name, category_description) VALUES (?, ?)")
	if err != nil {
		return false, "Internal server error"
	}
	res, err := statement.Exec(category.Name, category.Description)
	if err != nil {
		return false, err.Error()
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, "Internal server error"
	}
	if rowsAffected > 0 {
		return true, ""
	}
	return false, "Internal server error"
}

//Delete deletes category from the database
func (um *CategoryModel) Delete(id int) bool {
	res, err := um.DB.Exec("DELETE FROM categories WHERE category_id = ?", id)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}

//Update updates existing category in the database
func (um *CategoryModel) Update(category *entities.Category) bool {
	statement, err := um.DB.Prepare("UPDATE categories SET category_name = ?, category_description = ? WHERE category_id = ?")
	if err != nil {
		return false
	}
	res, err := statement.Exec(category.Name, category.Description, category.ID)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}
