package crud

import (
	"database/sql"
	"errors"
	"forum/api/models"
)

//CategoryModel helps performing CRUD operations
type CategoryModel struct {
	DB *sql.DB
}

//NewCategoryModel creates an instance of CategoryModel
func NewCategoryModel(db *sql.DB) *CategoryModel {
	return &CategoryModel{db}
}

//FindAll returns all categories in the database
func (um *CategoryModel) FindAll(postID uint64) ([]models.Category, error) {
	var (
		rows       *sql.Rows
		categories []models.Category
		err        error
	)
	if rows, err = um.DB.Query("SELECT * FROM categories WHERE category_post = ?", postID); err != nil {
		return nil, err
	}
	for rows.Next() {
		var category models.Category
		rows.Scan(&category.ID, &category.Post, &category.Name)
		categories = append(categories, category)
	}
	return categories, nil
}

//Find returns a specific category from the database
func (um *CategoryModel) Find(id int) (*models.Category, error) {
	var category models.Category
	rows, err := um.DB.Query("SELECT * FROM categories WHERE category_id = ?", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&category.ID, &category.Name)
	}
	return &category, nil
}

//Create adds a new category to the database
func (um *CategoryModel) Create(postID int64, categories ...string) (categoriesCreated int, err error) {
	var (
		stmt         *sql.Stmt
		result       sql.Result
		rowsAffected int64
	)
	if stmt, err = um.DB.Prepare("INSERT INTO categories (category_post, category_name) VALUES (?, ?)"); err != nil {
		return
	}
	for _, category := range categories {
		if result, err = stmt.Exec(postID, category); err != nil {
			return
		}
		categoriesCreated++
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return
	}
	if rowsAffected > 0 {
		return
	}
	return categoriesCreated, nil
}

//Update updates existing category in the database
func (um *CategoryModel) Update(category *models.Category) error {
	statement, err := um.DB.Prepare("UPDATE categories SET category_name = ?, category_description = ? WHERE category_id = ?")
	if err != nil {
		return err
	}
	res, err := statement.Exec(category.Name, category.ID)
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
func (um *CategoryModel) Delete(cid int) error {
	res, err := um.DB.Exec("DELETE FROM categories WHERE category_id = ?", cid)
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

func (cm *CategoryModel) DeleteGroup(pid uint64) error {
	var (
		res          sql.Result
		rowsAffected int64
		err          error
	)
	if res, err = cm.DB.Exec("DELETE FROM categories WHERE category_post = ?", pid); err != nil {
		return err
	}
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected > 0 {
		return nil
	}
	return errors.New("failed to delete replies for the post")
}
