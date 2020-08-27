package models

import (
	"database/sql"
	"errors"
	"forum/config"
	"net/http"
	"time"
)

//Post struct contains info about post
type Post struct {
	ID       int64
	By       uint64
	Category int
	Name     string
	Content  string
	Created  time.Time
	Updated  time.Time
	Rating   int
}

//PostModel helps performing CRUD operations
type PostModel struct {
	DB *sql.DB
}

//NewPostModel creates an instance of PostModel
func NewPostModel(db *sql.DB) (*PostModel, error) {
	return &PostModel{db}, nil
}

//FindAll returns all posts in the database
func (pm *PostModel) FindAll() ([]Post, error) {
	rows, err := pm.DB.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	var posts []Post

	for rows.Next() {
		var post Post
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
func (pm *PostModel) FindByID(uid int64) (*Post, int, error) {
	var post Post
	var created, lastOnline string
	row := pm.DB.QueryRow("SELECT * FROM posts WHERE post_id = ?", uid)
	err := row.Scan(&post.ID, &post.By, &post.Category, &post.Name, &post.Content, &post.Created, &post.Updated, &post.Rating)
	if err == sql.ErrNoRows {
		return nil, http.StatusBadRequest, errors.New("post not found")
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if post.Created, err = time.Parse(config.TimeLayout, created); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if post.Updated, err = time.Parse(config.TimeLayout, lastOnline); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &post, http.StatusOK, nil
}

//Create adds a new post to the database
func (pm *PostModel) Create(post *Post) error {
	statement, err := pm.DB.Prepare("INSERT INTO posts (post_by, post_category, post_name, post_content, post_created, post_updated, post_rating) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	res, err := statement.Exec(post.By, post.Category, post.Name, post.Content, time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout), post.Rating)
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
	return errors.New("could not create new post")
}

//Update updates existing post in the database
func (pm *PostModel) Update(post *Post) error {
	statement, err := pm.DB.Prepare("UPDATE posts SET post_by = ?, post_category = ?, post_name = ?, post_content = ?, post_created = ?, post_updated = ?, post_rating = ? = ? WHERE post_id = ?")
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
func (pm *PostModel) Delete(id int64) (int, error) {
	var err error
	var result sql.Result
	var rowsAffected int64
	result, err = pm.DB.Exec("DELETE FROM posts WHERE post_id = ?", id)
	if err == sql.ErrNoRows {
		return http.StatusBadRequest, errors.New("user not found")
	}
	if err != nil {
		return http.StatusInternalServerError, err
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
