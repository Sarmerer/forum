package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.PostRepo = crud.NewPostRepoCRUD()
		userCtx models.UserCtx      = utils.GetUserFromCtx(r)
		posts   []models.Post
		err     error
	)
	if posts, err = repo.FindAll(userCtx.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, nil, posts)
}

func FindPost(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.PostRepo = crud.NewPostRepoCRUD()
		userCtx models.UserCtx      = utils.GetUserFromCtx(r)
		input   models.InputPostFind
		posts   interface{}
		status  int
		err     error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	switch input.By {
	case "id":
		var post *models.Post
		if post, status, err = repo.FindByID(input.ID, userCtx.ID); err != nil {
			response.Error(w, status, err)
			return
		}
		response.Success(w, nil, post)
		return
	case "author":
		if posts, err = repo.FindByAuthor(input.AuthorID); err != nil {
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
		repo    repository.PostRepo = crud.NewPostRepoCRUD()
		userCtx models.UserCtx      = utils.GetUserFromCtx(r)
		input   models.InputPostCreateUpdate
		post    models.Post
		pid     int64
		err     error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	post = models.Post{
		Title:      input.Title,
		Content:    input.Content,
		AuthorID:   userCtx.ID,
		AuthorName: userCtx.DisplayName,
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
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		input  models.InputPostCreateUpdate
		pid    int64
		post   *models.Post
		status int
		err    error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if pid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if post, status, err = repo.FindByID(pid, -1); err != nil {
		response.Error(w, status, err)
		return
	}

	if !requestorIsEntityOwner(utils.GetUserFromCtx(r), post.AuthorID) {
		response.Error(w, http.StatusForbidden, errors.New("this post doesn't belong to you"))
		return
	}

	post.Updated = time.Now().Format(config.TimeLayout)
	if input.Title != "" {
		post.Title = input.Title
	}
	if input.Content != "" {
		post.Content = input.Content
	}
	if err = repo.Update(post); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, fmt.Sprint("post has been updated"), nil)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		pid    int64
		post   *models.Post
		status int
		err    error
	)
	if pid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if post, status, err = repo.FindByID(pid, -1); err != nil {
		response.Error(w, status, err)
		return
	}

	if !requestorIsEntityOwner(utils.GetUserFromCtx(r), post.AuthorID) {
		response.Error(w, http.StatusForbidden, errors.New("this post doesn't belong to you"))
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
