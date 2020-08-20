package models

import (
	"database/sql"
	"forum/api/entities"
	"forum/config"
	"time"
)

//PostModel helps performing CRUD operations
type PostModel struct {
	DB *sql.DB
}

//NewPostModel creates an instance of PostModel
func NewPostModel(db *sql.DB) (*PostModel, error) {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "posts" (
	"post_id"	INTEGER,
	"post_by"	INTEGER,
	"post_category"	INTEGER,
	"post_name"	TEXT,
	"post_content"	TEXT,
	"post_created"	TEXT,
	"post_updated"	TEXT,
	"post_likes"	INTEGER,
	"post_dislikes"	INTEGER,
	FOREIGN KEY("post_by") REFERENCES "users"("user_id"),
	FOREIGN KEY("post_category") REFERENCES "categories"("category_id"),
	PRIMARY KEY("post_id" AUTOINCREMENT)
)`)
	if err != nil {
		return nil, err
	}
	statement.Exec()
	return &PostModel{
		DB: db,
	}, nil
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
		rows.Scan(&post.ID, &post.By, &post.Category, &post.Name, &post.Content, &created, &updated, &post.Likes, &post.Dislikes)
		date, _ := time.Parse(config.TimeLayout, created)
		post.Created = date
		date, _ = time.Parse(config.TimeLayout, updated)
		posts = append(posts, post)
	}
	return posts, nil
}

//Find returns a specific post from the database
func (um *PostModel) Find(id int) (entities.Post, error) {
	var post entities.Post
	rows, err := um.DB.Query("SELECT * FROM posts WHERE post_id = ?", id)
	if err != nil {
		return post, err
	}
	for rows.Next() {
		var created, updated string
		rows.Scan(&post.ID, &post.By, &post.Category, &post.Name, &post.Content, &created, &updated, &post.Likes, &post.Dislikes)
		date, _ := time.Parse(config.TimeLayout, created)
		post.Created = date
		date, _ = time.Parse(config.TimeLayout, updated)
		post.Updated = date
	}
	return post, nil
}

//Create adds a new post to the database
func (um *PostModel) Create(post *entities.Post) (bool, string) {
	statement, err := um.DB.Prepare("INSERT INTO posts (post_by, post_category, post_name, post_content, post_created, post_updated, post_likes, post_dislikes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return false, "Internal server error"
	}
	res, err := statement.Exec(post.By, post.Category, post.Name, post.Content, time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout), post.Likes, post.Dislikes)
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

//Delete deletes post from the database
func (um *PostModel) Delete(id int) bool {
	res, err := um.DB.Exec("DELETE FROM posts WHERE post_id = ?", id)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}

//Update updates existing post in the database
func (um *PostModel) Update(post *entities.Post) bool {
	statement, err := um.DB.Prepare("UPDATE posts SET post_by = ?, post_category = ?, post_name = ?, post_content = ?, post_created = ?, post_updated = ?, post_likes = ?, post_dislikes = ? WHERE post_id = ?")
	if err != nil {
		return false
	}
	res, err := statement.Exec(post.By, post.Category, post.Name, post.Content, post.Created.Format(config.TimeLayout), post.Updated.Format(config.TimeLayout), post.Likes, post.Dislikes, post.ID)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}
