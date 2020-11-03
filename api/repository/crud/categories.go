package crud

import (
	"database/sql"
	"errors"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
)

//CategoryRepoCRUD helps performing CRUD operations
type CategoryRepoCRUD struct{}

//NewCategoryRepoCRUD creates an instance of CategoryModel
//FIXME wrong categories amount after COUNT
func NewCategoryRepoCRUD() CategoryRepoCRUD {
	return CategoryRepoCRUD{}
}

func (CategoryRepoCRUD) FindAll() ([]models.Category, error) {
	var (
		rows       *sql.Rows
		categories []models.Category
		err        error
	)
	if rows, err = repository.DB.Query(
		`SELECT COUNT(category_id_fkey) AS use_count,
			category_id_fkey,
			name
		FROM categories c
			LEFT JOIN posts_categories_bridge ctb ON ctb.category_id_fkey = c.id
		GROUP BY category_id_fkey
		ORDER BY name ASC
		LIMIT 30`,
	); err != nil {
		return nil, err
	}
	for rows.Next() {
		var c models.Category
		rows.Scan(&c.UseCount, &c.ID, &c.Name)
		categories = append(categories, c)
	}
	return categories, nil
}

//FindByPostID returns all categories belonging to  a post
func (CategoryRepoCRUD) FindByPostID(postID int64) ([]models.Category, error) {
	var (
		rows       *sql.Rows
		categories []models.Category
		err        error
	)
	if rows, err = repository.DB.Query(
		`SELECT category_id_fkey,
			name
		FROM categories c
			LEFT JOIN posts_categories_bridge ctb ON ctb.post_id_fkey = ?
		WHERE c.id = ctb.category_id_fkey`,
		postID,
	); err != nil {
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
func (CategoryRepoCRUD) Find(id int) (*models.Category, error) {
	var category models.Category
	rows, err := repository.DB.Query("SELECT * FROM categories WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&category.ID, &category.Name)
	}
	return &category, nil
}

//Create adds a new category to the database
//FIXME category being duplicated, when creating new post with existing category
func (CategoryRepoCRUD) Create(postID int64, categories []string) (err error) {
	var (
		cid    int64
		result sql.Result
	)
	for _, category := range categories {
		if err = repository.DB.QueryRow(
			`SELECT id
			FROM categories
			WHERE name = ?`,
			category,
		).Scan(
			&cid,
		); err != nil {
			if err != sql.ErrNoRows {
				return err
			}
			if result, err = repository.DB.Exec(
				`INSERT INTO categories (name)
				VALUES (?)`,
				category,
			); err != nil {
				return err
			}
			if cid, err = result.LastInsertId(); err != nil {
				return err
			}
		}
		if _, err = repository.DB.Exec(
			`INSERT INTO posts_categories_bridge (post_id_fkey, category_id_fkey)
			VALUES (?, ?)`,
			postID, cid,
		); err != nil {
			return err
		}
	}
	return nil
}

//Delete deletes category from the database
//TODO also delete category from the brige table
// FIXME BEGIN COMMIT do not work for some reason
func (CategoryRepoCRUD) Delete(cid int) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`BEGIN;
		DELETE FROM categories
		WHERE id = $1;
		DELETE FROM posts_categories_bridge
		WHERE category_id_fkey = $1;
		COMMIT;`,
		cid,
	); err != nil {
		return err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("could not delete the category")
}

func (CategoryRepoCRUD) DeleteGroup(pid int64) error {
	var err error
	if _, err = repository.DB.Exec(
		`DELETE FROM posts_categories_bridge
		WHERE post_id_fkey = ?`,
		pid,
	); err != nil {
		return err
	}
	if _, err = repository.DB.Exec(
		`DELETE FROM categories
		WHERE id IN (
				SELECT c.id
				FROM categories c
					LEFT JOIN posts_categories_bridge pcb ON c.id = pcb.category_id_fkey
				WHERE pcb.category_id_fkey IS NULL
			)`,
	); err != nil {
		return err
	}
	return nil
}
