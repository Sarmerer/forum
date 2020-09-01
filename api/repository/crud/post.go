package crud

import (
	"database/sql"
	"errors"
	"forum/api/models"
	"forum/config"
	"net/http"
	"time"
)

//PostModel helps performing CRUD operations
type PostModel struct {
	DB *sql.DB
}

//NewPostModel creates an instance of PostModel
func NewPostModel(db *sql.DB) *PostModel {
	return &PostModel{db}
}

//FindAll returns all posts in the database
func (pm *PostModel) FindAll() ([]models.Post, error) {
	rows, err := pm.DB.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	var posts []models.Post

	for rows.Next() {
		var post models.Post
		rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.Created, &post.Updated, &post.Rating)
		posts = append(posts, post)
	}
	return posts, nil
}

//FindByID returns a specific post from the database
func (pm *PostModel) FindByID(pid uint64) (*models.Post, int, error) {
	var (
		post models.Post
		err  error
	)
	if err = pm.DB.
		QueryRow("SELECT * FROM posts WHERE id = ?", pid).
		Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.Created, &post.Updated, &post.Rating); err == sql.ErrNoRows {
		return nil, http.StatusBadRequest, errors.New("post not found")
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &post, http.StatusOK, nil
}

//Create adds a new post to the database
func (pm *PostModel) Create(post *models.Post) (int64, error) {
	var (
		result       sql.Result
		rowsAffected int64
		pid          int64
		err          error
	)
	if result, err = pm.DB.Exec(
		"INSERT INTO posts (author_fkey, title, content, created, updated, rating) VALUES (?, ?, ?, ?, ?, ?)",
		post.Author, post.Title, post.Content, time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout), post.Rating); err != nil {
		return 0, err
	}
	if pid, err = result.LastInsertId(); err != nil {
		return 0, err
	}
	if rowsAffected, err = result.RowsAffected(); rowsAffected > 0 && err == nil {
		return pid, nil
	}
	return 0, err
}

//Update updates existing post in the database
func (pm *PostModel) Update(post *models.Post) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = pm.DB.Exec(
		"UPDATE posts SET author_fkey = ?, title = ?, content = ?, created = ?, updated = ?, rating = ? WHERE id = ?",
		post.Author, post.Title, post.Content, post.Created, post.Updated, post.Rating, post.ID); err != nil {
		return err
	}
	if rowsAffected, err = result.RowsAffected(); rowsAffected > 0 && err == nil {
		return nil
	}
	return err
}

//Delete deletes post from the database
func (pm *PostModel) Delete(pid uint64) (int, error) {
	var err error
	var result sql.Result
	var rowsAffected int64
	result, err = pm.DB.Exec("DELETE FROM posts WHERE id = ?", pid)
	if err == sql.ErrNoRows {
		return http.StatusBadRequest, errors.New("post not found")
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); rowsAffected > 0 && err == nil {
		return http.StatusOK, nil
	}
	return http.StatusInternalServerError, err
}
