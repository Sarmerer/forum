package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/api/config"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/repository/crud"
	"forum/api/response"
	"forum/api/utils"
	"net/http"
	"time"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var (
		repo  repository.PostRepo = crud.NewPostRepoCRUD()
		uid   int64               = utils.GetUIDFromCtx(r)
		posts []models.Post
		err   error
	)
	if posts, err = repo.FindAll(uid); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, nil, posts)
}

func FindPost(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		uid    int64               = utils.GetUIDFromCtx(r)
		input  models.InputPostFind
		posts  interface{}
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
		if post, status, err = repo.FindByID(input.ID, uid); err != nil {
			response.Error(w, status, err)
			return
		}
		response.Success(w, nil, post)
		return
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
	response.Success(w, nil, posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		uid    int64               = r.Context().Value("uid").(int64)
		author *models.User
		input  models.InputPostCreateUpdate
		post   models.Post
		pid    int64
		status int
		err    error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if author, status, err = crud.NewUserRepoCRUD().FindByID(uid); err != nil {
		response.Error(w, status, err)
		return
	}
	post = models.Post{
		Title:      input.Title,
		Content:    input.Content,
		AuthorID:   uid,
		AuthorName: author.DisplayName,
		Created:    time.Now().Format(config.TimeLayout),
		Updated:    time.Now().Format(config.TimeLayout),
		Rating:     0,
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
		input       models.InputPostCreateUpdate
		name        string
		content     string
		pid         int64
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
	if updatedPost, status, err = repo.FindByID(pid, -1); err != nil {
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
		pid    int64
		status int
		err    error
	)
	if pid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if _, status, err = repo.FindByID(pid, -1); err != nil {
		response.Error(w, status, err)
		return
	}
	if err = DeleteCommentsGroup(pid); err != nil {
		fmt.Println(err)
	}
	//TODO make proper logging here
	if err = DeleteAllCategoriesForPost(pid); err != nil {
		fmt.Println(err)
	}
	if err = DeleteReactionsForPost(pid); err != nil {
		fmt.Println(err)
	}
	if status, err = repo.Delete(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, fmt.Sprintf("post has been deleted"), nil)
}
