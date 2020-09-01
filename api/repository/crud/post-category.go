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
	if rows, err = um.DB.Query("SELECT category_id_fkey, name FROM categories c LEFT JOIN posts_categories_bridge ctb ON ctb.post_id_fkey = ? WHERE c.id = ctb.category_id_fkey", postID); err != nil {
		return nil, err
	}
	for rows.Next() {
		var category models.Category
		rows.Scan(&category.ID, &category.Name)
		categories = append(categories, category)
	}
	return categories, nil
}

//Find returns a specific category from the database
//TODO implment search for all post with such categories here
func (um *CategoryModel) Find(id int) (*models.Category, error) {
	var category models.Category
	rows, err := um.DB.Query("SELECT * FROM categories WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&category.ID, &category.Name)
	}
	return &category, nil
}

//Create adds a new category to the database
func (cm *CategoryModel) Create(postID int64, categories ...string) (categoriesCreated int, err error) {
	var (
		cid    int64
		result sql.Result
	)
	for _, category := range categories {
		if err = cm.DB.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&cid); err != nil {
			if err != sql.ErrNoRows {
				return 0, err
			}
			if result, err = cm.DB.Exec("INSERT INTO categories (name) VALUES (?)", category); err != nil {
				return 0, err
			}
			if cid, err = result.LastInsertId(); err != nil {
				return 0, err
			}
		}
		if _, err = cm.DB.Exec("INSERT INTO posts_categories_bridge (post_id_fkey, category_id_fkey) VALUES (?, ?)", postID, cid); err != nil {
			return 0, err
		}
		categoriesCreated++
	}
	return categoriesCreated, nil
}

//Delete deletes category from the database
func (um *CategoryModel) Delete(cid int) error {
	res, err := um.DB.Exec("DELETE FROM categories WHERE id = ?", cid)
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
		err error
	)
	if _, err = cm.DB.Exec("DELETE FROM categories WHERE post_id_fkey = ?", pid); err != nil {
		return err
	}
	if _, err = cm.DB.Exec("DELETE FROM posts_categories_bridge WHERE post_id_fkey = ?", pid); err != nil {
		return err
	}
	return nil
}
