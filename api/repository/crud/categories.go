package crud

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
)

//CategoryRepoCRUD helps performing CRUD operations
type CategoryRepoCRUD struct{}

//NewCategoryRepoCRUD creates an instance of CategoryModel
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
			LEFT JOIN posts_categories_bridge ctb ON ctb.category_id_fkey = c._id
		GROUP BY category_id_fkey
		ORDER BY use_count DESC,name ASC`,
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
func (CategoryRepoCRUD) FindByPostID(postID int64) ([]models.Category, int, error) {
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
		WHERE c._id = ctb.category_id_fkey`,
		postID,
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	for rows.Next() {
		var category models.Category
		rows.Scan(&category.ID, &category.Name)
		categories = append(categories, category)
	}
	return categories, http.StatusOK, nil
}

//Find returns a specific category from the database
// ? do we even need that
func (CategoryRepoCRUD) Find(id int) (*models.Category, error) {
	var category models.Category
	rows, err := repository.DB.Query("SELECT * FROM categories WHERE _id = ?", id)
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
			`SELECT _id
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

func (CategoryRepoCRUD) Update(postID int64, categories []string) error {
	var (
		repo CategoryRepoCRUD = NewCategoryRepoCRUD()
		err  error
	)
	if err = repo.DeleteGroup(postID); err != nil {
		return err
	}
	if err = repo.Create(postID, categories); err != nil {
		return err
	}
	return nil
}

//Delete deletes category from the database
func (CategoryRepoCRUD) Delete(categotyID int) error {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
	)
	ctx = context.Background()
	tx, err = repository.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		`DELETE FROM categories
		WHERE _id = ?`, categotyID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx,
		`DELETE FROM posts_categories_bridge
		WHERE category_id_fkey = ?`, categotyID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
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
		WHERE _id IN (
				SELECT c._id
				FROM categories c
					LEFT JOIN posts_categories_bridge pcb ON c._id = pcb.category_id_fkey
				WHERE pcb.category_id_fkey IS NULL
			)`,
	); err != nil {
		return err
	}
	return nil
}
