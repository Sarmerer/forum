package crud

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/api/models"
	"forum/api/repository"
	"forum/config"
	"net/http"
	"time"
)

//PostRepoCRUD helps performing CRUD operations
type PostRepoCRUD struct{}

//NewPostRepoCRUD creates an instance of PostModel
func NewPostRepoCRUD() PostRepoCRUD {
	return PostRepoCRUD{}
}

//FindAll returns all posts in the database
func (PostRepoCRUD) FindAll() ([]models.Post, error) {
	var (
		rows  *sql.Rows
		posts []models.Post
		err   error
	)
	if rows, err = repository.DB.Query("SELECT * FROM posts"); err != nil {
		return nil, err
	}
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.AuthorName, &p.Title, &p.Content, &p.Created, &p.Updated, &p.Rating)
		posts = append(posts, p)
	}
	return posts, nil
}

//FindByID returns a specific post from the database
func (PostRepoCRUD) FindByID(pid uint64) (*models.Post, int, error) {
	var (
		p   models.Post
		err error
	)
	if err = repository.DB.QueryRow(
		"SELECT * FROM posts WHERE id = ?", pid,
	).Scan(
		&p.ID, &p.AuthorID, &p.AuthorName, &p.Title, &p.Content, &p.Created, &p.Updated, &p.Rating,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, errors.New("post not found")
	}
	return &p, http.StatusOK, nil
}

func (PostRepoCRUD) FindByAuthor(uid uint64) ([]models.Post, error) {
	var (
		rows  *sql.Rows
		posts []models.Post
		err   error
	)
	if rows, err = repository.DB.Query(
		"SELECT * FROM posts WHERE author_fkey = ?",
		uid,
	); err != nil {
		return nil, err
	}
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.AuthorName, &p.Title, &p.Content, &p.Created, &p.Updated, &p.Rating)
		posts = append(posts, p)
	}
	return posts, nil
}

func (PostRepoCRUD) FindByCategories(categories []string) ([]models.Post, error) {
	var (
		rows  *sql.Rows
		args  string
		posts []models.Post
		err   error
	)
	for _, category := range categories {
		if len(args) == 0 {
			args += fmt.Sprintf(`"%s"`, category)
		} else {
			args += fmt.Sprintf(`, "%s"`, category)
		}
	}
	if rows, err = repository.DB.Query(
		`SELECT p.id, p.author_fkey, p.author_name_fkey, p.title, p.content, p.created, p.updated, p.rating
		FROM posts_categories_bridge AS pcb
		INNER JOIN posts as p
			ON p.id = pcb.post_id_fkey
		INNER JOIN categories AS c
			ON c.id = pcb.category_id_fkey
			   WHERE c.name IN (`+args+`)
		GROUP BY
			p.id
		HAVING 
			COUNT(DISTINCT c.id) = ?`,
		len(categories),
	); err != nil {
		return nil, err
	}
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, p.AuthorName, &p.Title, &p.Content, &p.Created, &p.Updated, &p.Rating)
		posts = append(posts, p)
	}
	return posts, nil
}

//Create adds a new post to the database
func (PostRepoCRUD) Create(post *models.Post) (int64, error) {
	var (
		result       sql.Result
		rowsAffected int64
		pid          int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"INSERT INTO posts (author_fkey,author_name_fkey, title, content, created, updated, rating) VALUES (?, ?, ?, ?, ?, ?, ?)",
		post.AuthorID, post.AuthorName, post.Title, post.Content, time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout), post.Rating,
	); err != nil {
		return 0, err
	}

	if pid, err = result.LastInsertId(); err != nil {
		return 0, err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return 0, err
	}
	if rowsAffected > 0 {
		return pid, nil
	}
	return 0, errors.New("could not create the post")
}

//Update updates existing post in the database
func (PostRepoCRUD) Update(post *models.Post) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"UPDATE posts SET author_fkey = ?, author_name_fkey, title = ?, content = ?, created = ?, updated = ?, rating = ? WHERE id = ?",
		post.AuthorID, post.AuthorName, post.Title, post.Content, post.Created, post.Updated, post.Rating, post.ID,
	); err != nil {
		return err
	}

	if rowsAffected, err = result.RowsAffected(); err == nil {
		return nil
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("could not update the post")
}

//Delete deletes post from the database
func (PostRepoCRUD) Delete(pid uint64) (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		"DELETE FROM posts WHERE id = ?", pid,
	); err != nil {
		if err != sql.ErrNoRows {
			return http.StatusInternalServerError, err
		}
		return http.StatusNotFound, errors.New("post not found")
	}

	if rowsAffected, err = result.RowsAffected(); err == nil {
		return http.StatusOK, nil
	}
	if rowsAffected > 0 {
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not delete the post")
}
