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
		var created, updated string
		rows.Scan(&post.ID, &post.By, &post.Category, &post.Name, &post.Content, &created, &updated, &post.Rating)
		date, _ := time.Parse(config.TimeLayout, created)
		post.Created = date
		date, _ = time.Parse(config.TimeLayout, updated)
		posts = append(posts, post)
	}
	return posts, nil
}

//Find returns a specific post from the database
func (pm *PostModel) FindByID(pid uint64) (*models.Post, int, error) {
	var post models.Post
	var created, updated string
	row := pm.DB.QueryRow("SELECT * FROM posts WHERE post_id = ?", pid)
	err := row.Scan(&post.ID, &post.By, &post.Category, &post.Name, &post.Content, &created, &updated, &post.Rating)
	if err == sql.ErrNoRows {
		return nil, http.StatusBadRequest, errors.New("post not found")
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if post.Created, err = time.Parse(config.TimeLayout, created); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if post.Updated, err = time.Parse(config.TimeLayout, updated); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &post, http.StatusOK, nil
}

//Create adds a new post to the database
func (pm *PostModel) Create(post *models.Post) (pid int64, err error) {
	var (
		stmt         *sql.Stmt
		result       sql.Result
		rowsAffected int64
	)
	if stmt, err = pm.DB.Prepare("INSERT INTO posts (post_by, post_category, post_name, post_content, post_created, post_updated, post_rating) VALUES (?, ?, ?, ?, ?, ?, ?)"); err != nil {
		return
	}
	if result, err = stmt.Exec(post.By, post.Category, post.Name, post.Content, time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout), post.Rating); err != nil {
		return
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return
	}
	if pid, err = result.LastInsertId(); err != nil {
		return
	}
	if rowsAffected > 0 {
		return
	}
	return 0, errors.New("could not create new post")
}

//Update updates existing post in the database
func (pm *PostModel) Update(post *models.Post) error {
	statement, err := pm.DB.Prepare("UPDATE posts SET post_by = ?, post_category = ?, post_name = ?, post_content = ?, post_created = ?, post_updated = ?, post_rating = ? WHERE post_id = ?")
	if err != nil {
		return err
	}
	res, err := statement.Exec(post.By, post.Category, post.Name, post.Content, post.Created.Format(config.TimeLayout), post.Updated.Format(config.TimeLayout), post.Rating, post.ID)
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
	return errors.New("could not update the post")
}

//Delete deletes post from the database
func (pm *PostModel) Delete(pid uint64) (int, error) {
	var err error
	var result sql.Result
	var rowsAffected int64
	result, err = pm.DB.Exec("DELETE FROM posts WHERE post_id = ?", pid)
	if err == sql.ErrNoRows {
		return http.StatusBadRequest, errors.New("post not found")
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); rowsAffected > 0 && err == nil {
		return http.StatusOK, nil
	} else {
		return http.StatusInternalServerError, err
	}
}
