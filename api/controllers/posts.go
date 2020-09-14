package controllers

import (
	"encoding/json"
	"errors"
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

type postResponse struct {
	Post       *models.Post `json:"post"`
	Categories interface{}  `json:"categories"`
	Replies    interface{}  `json:"replies"`
}

type createUpdateInput struct {
	Title      string
	Content    string
	Categories []string
}

type findInput struct {
	By         string   `json:"by"`
	ID         uint64   `json:"id"`
	Author     uint64   `json:"author"`
	Categories []string `json:"categories"`
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		posts  []models.Post
		result []postResponse
		err    error
	)
	if posts, err = repo.FindAll(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	for i := 0; i < len(posts); i++ {
		p := postResponse{Post: &posts[i]}
		if p.Categories, err = GetCategoriesByPostID(posts[i].ID); err != nil {
			p.Categories = err
		}
		if p.Replies, err = CountReplies(posts[i].ID); err != nil {
			p.Replies = err
		}
		result = append(result, p)
	}
	response.Success(w, nil, result)
}

func FindPost(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		input  findInput
		result []postResponse
		posts  []models.Post
		status int
		err    error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	switch input.By {
	case "id":
		var post *models.Post
		if post, status, err = repo.FindByID(input.ID); err != nil {
			response.Error(w, status, err)
			return
		}
		posts = append(posts, *post)
	case "author":
		if posts, err = repo.FindByAuthor(input.Author); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
	case "categories":
		if posts, err = repo.FindByCategories(input.Categories); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
	default:
		response.Error(w, http.StatusBadRequest, errors.New("unknown search type"))
		return
	}
	for i := 0; i < len(posts); i++ {
		p := postResponse{Post: &posts[i]}
		if p.Categories, err = GetCategoriesByPostID(posts[i].ID); err != nil {
			p.Categories = err
		}
		if p.Replies, err = CountReplies(posts[i].ID); err != nil {
			p.Replies = err
		}
		result = append(result, p)
	}
	response.Success(w, nil, result)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		author uint64              = r.Context().Value("uid").(uint64)
		input  createUpdateInput
		post   models.Post
		pid    int64
		err    error
	)
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
		if err = crud.NewCategoryRepoCRUD().Create(pid, input.Categories); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
	}

	response.Success(w, fmt.Sprintf("post has been created"), pid)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo        repository.PostRepo = crud.NewPostRepoCRUD()
		input       createUpdateInput
		name        string
		content     string
		pid         uint64
		updatedPost *models.Post
		status      int
		err         error
	)
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
	if err = DeleteAllRepliesForPost(pid); err != nil {
		fmt.Println(err)
	}
	if err = DeleteAllCategoriesForPost(pid); err != nil {
		fmt.Println(err)
	}
	if status, err = repo.Delete(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, fmt.Sprintf("post has been deleted"), nil)
}
