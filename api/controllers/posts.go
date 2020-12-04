package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/logger"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.PostRepo = crud.NewPostRepoCRUD()
		input   models.InputAllPosts
		userCtx models.UserCtx = utils.GetUserFromCtx(r)
		result  *models.Posts
		err     error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	input.Validate()

	if result, err = repo.FindAll(userCtx.ID, input); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, nil, result)
}

func FindPost(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.PostRepo = crud.NewPostRepoCRUD()
		userCtx models.UserCtx      = utils.GetUserFromCtx(r)
		input   models.InputFind
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
		if posts, status, err = repo.FindByAuthor(input.AuthorID, userCtx.ID); err != nil {
			response.Error(w, status, err)
			return
		}
	case "categories":
		if posts, err = repo.FindByCategories(input.Categories, userCtx.ID); err != nil {
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
		newPost *models.Post
		status  int
		err     error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	post = models.Post{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: userCtx.ID,
		Created:  time.Now().Format(config.TimeLayout),
		Updated:  time.Now().Format(config.TimeLayout),
		Rating:   0,
	}
	if newPost, status, err = repo.Create(&post, input.Categories); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, fmt.Sprintf("post has been created"), newPost)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.PostRepo = crud.NewPostRepoCRUD()
		input   models.InputPostCreateUpdate
		post    *models.Post
		userCtx models.UserCtx = utils.GetUserFromCtx(r)
		status  int
		err     error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if post, status, err = repo.FindByID(input.ID, -1); err != nil {
		response.Error(w, status, err)
		return
	}

	if !requestorIsEntityOwner(utils.GetUserFromCtx(r), post.AuthorID) {
		response.Error(w, http.StatusForbidden, errors.New("this post doesn't belong to you"))
		return
	}

	post.Title = input.Title
	post.Content = input.Content

	if err = crud.NewCategoryRepoCRUD().Update(input.ID, input.Categories); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if post, status, err = repo.Update(post, userCtx); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, fmt.Sprint("post has been updated"), post)
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
		logger.CheckErrAndLog("comments deletion", "", err)
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if err = DeleteAllCategoriesForPost(pid); err != nil {
		logger.CheckErrAndLog("categories deletion", "", err)
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if err = DeleteReactionsForPost(pid); err != nil {
		logger.CheckErrAndLog("reactions deletion", "", err)
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if status, err = repo.Delete(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, fmt.Sprintf("post has been deleted"), nil)
}
