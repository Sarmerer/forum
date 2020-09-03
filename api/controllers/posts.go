package controllers

import (
	"encoding/json"
	"fmt"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/api/utils"
	"forum/config"
	"net/http"
	"time"
)

type responsePost struct {
	Post       *models.Post `json:"post"`
	Categories interface{}  `json:"categories"`
	Replies    interface{}  `json:"replies"`
}

//TODO improve response time
func GetPosts(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		posts  []models.Post
		result []responsePost
		err    error
	)
	if posts, err = repo.FindAll(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	for i := 0; i < len(posts); i++ {
		p := responsePost{Post: &posts[i]}
		if p.Replies, err = CountReplies(posts[i].ID); err != nil {
			p.Replies = err
		}
		if p.Categories, err = GetCategories(posts[i].ID); err != nil {
			p.Categories = err
		}
		result = append(result, p)
	}
	response.Success(w, nil, result)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		pid    uint64
		result responsePost
		status int
		err    error
	)
	if pid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if result.Post, status, err = repo.FindByID(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	if result.Categories, err = GetCategories(pid); err != nil {
		result.Categories = err
	}
	if result.Replies, err = GetReplies(pid); err != nil {
		result.Replies = err
	}
	response.Success(w, nil, result)
}

func FindPost(w http.ResponseWriter, r *http.Request) {
	var (
		repo  repository.PostRepo = crud.NewPostRepoCRUD()
		posts []models.Post
		err   error
	)
	input := struct {
		Categories []string `json:"categories"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
	}
	if posts, err = repo.FindByCategories(input.Categories); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, nil, posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		author uint64              = r.Context().Value("uid").(uint64)
		post   models.Post
		pid    int64
		err    error
	)
	input := struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Categories []string `json:"categories"`
	}{}

	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	post = models.Post{
		Title:   input.Title,
		Content: input.Content,
		Author:  author,
		Created: time.Now().Format(config.TimeLayout),
		Updated: time.Now().Format(config.TimeLayout),
		Rating:  0,
	}
	if pid, err = repo.Create(&post); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if len(input.Categories) > 0 {
		if err = crud.NewCategoryRepoCRUD().Create(pid, input.Categories...); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
	}

	response.Success(w, fmt.Sprintf("post has been created"), nil)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo        repository.PostRepo = crud.NewPostRepoCRUD()
		name        string
		content     string
		pid         uint64
		updatedPost *models.Post
		status      int
		err         error
	)
	input := struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Categories []string `json:"categories"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if pid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if updatedPost, status, err = repo.FindByID(pid); err != nil {
		response.Error(w, status, err)
		return
	}

	updatedPost.Updated = time.Now().Format(config.TimeLayout)
	if input.Title != "" {
		updatedPost.Title = name
	}
	if input.Content != "" {
		updatedPost.Content = content
	}
	if err = repo.Update(updatedPost); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, fmt.Sprint("post has been updated"), nil)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		pid    uint64
		status int
		err    error
	)
	if pid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if _, status, err = repo.FindByID(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	if status, err = repo.Delete(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	if err = DeleteAllRepliesForPost(pid); err != nil {
		fmt.Println(err)
	}
	if err = DeleteAllCategoriesForPost(pid); err != nil {
		fmt.Println(err)
	}
	response.Success(w, fmt.Sprintf("post has been deleted"), nil)
}
