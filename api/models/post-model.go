package models

import (
	"database/sql"
	"errors"
	"forum/api/entities"
	"forum/config"
	"forum/database"
	"time"
)

//PostModel helps performing CRUD operations
type PostModel struct {
	DB *sql.DB
}

//NewPostModel creates an instance of PostModel
func NewPostModel() (*PostModel, error) {
	db, dbErr := database.Connect()
	if dbErr != nil {
		return nil, dbErr
	}
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "posts" (
	"post_id"	INTEGER,
	"post_by"	INTEGER,
	"post_category"	INTEGER,
	"post_name"	TEXT,
	"post_content"	TEXT,
	"post_created"	TEXT,
	"post_updated"	TEXT,
	"post_rating"	INTEGER,
	FOREIGN KEY("post_by") REFERENCES "users"("user_id"),
	FOREIGN KEY("post_category") REFERENCES "categories"("category_id"),
	PRIMARY KEY("post_id" AUTOINCREMENT))`,
	)
	if err != nil {
		return nil, err
	}
	statement.Exec()
	return &PostModel{db}, nil
}

//FindAll returns all posts in the database
func (um *PostModel) FindAll() ([]entities.Post, error) {
	rows, e := um.DB.Query("SELECT * FROM posts")
	if e != nil {
		return nil, e
	}
	var posts []entities.Post

	for rows.Next() {
		var post entities.Post
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
func (um *PostModel) Find(id int64) (*entities.Post, error) {
	var post entities.Post
	rows, err := um.DB.Query("SELECT * FROM posts WHERE post_id = ?", id)
	if err != nil {
		return &post, err
	}
	for rows.Next() {
		var created, updated string
		rows.Scan(&post.ID, &post.By, &post.Category, &post.Name, &post.Content, &created, &updated, &post.Rating)
		date, _ := time.Parse(config.TimeLayout, created)
		post.Created = date
		date, _ = time.Parse(config.TimeLayout, updated)
		post.Updated = date
	}
	return &post, nil
}

//Create adds a new post to the database
func (um *PostModel) Create(post *entities.Post) error {
	statement, err := um.DB.Prepare("INSERT INTO posts (post_by, post_category, post_name, post_content, post_created, post_updated, post_rating) VALUES (?, ?, ?, ?, ?, ?, ?)")
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
func (um *PostModel) Update(post *entities.Post) error {
	statement, err := um.DB.Prepare("UPDATE posts SET post_by = ?, post_category = ?, post_name = ?, post_content = ?, post_created = ?, post_updated = ?, post_rating = ? = ? WHERE post_id = ?")
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
func (um *PostModel) Delete(id int64) error {
	res, err := um.DB.Exec("DELETE FROM posts WHERE post_id = ?", id)
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
	return errors.New("could not delete the post")
}
