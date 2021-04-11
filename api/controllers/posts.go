package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sarmerer/forum/api/logger"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"
)

// GetPosts returns struct with fields:
//
// - hot - several posts, sorted by their rating
//
// - recent - several posts sorted by date
//
// - total_rows - total amount of rows in posts table.
// Used for pagination component on frontend
func GetPosts(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.PostRepo = crud.NewPostRepoCRUD()
		input   models.InputAllPosts
		userCtx models.UserCtx = utils.GetUserFromCtx(r)
		result  *models.Posts
		status  int
		err     error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	input.Validate()

	if result, status, err = repo.FindAll(userCtx.ID, &input); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, nil, result)
}

// FindPost allows to search posts by various parameters.
// Currently supported:
//
// - id - returns a post with that id
//
// - author - returns all posts from single user
//
// - categories - returns posts that have theese categories
func FindPost(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.PostRepo = crud.NewPostRepoCRUD()
		userCtx models.UserCtx      = utils.GetUserFromCtx(r)
		input   models.InputFindPost
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
		if posts, status, err = repo.FindByCategories(input.Categories, userCtx.ID); err != nil {
			response.Error(w, status, err)
			return
		}
	default:
		response.Error(w, http.StatusBadRequest, errors.New("unknown search type"))
		return
	}
	response.Success(w, nil, posts)
}

// CreatePost adds new post record to database
//
// Returns Post model on success
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
		IsImage:  input.IsImage,
		AuthorID: userCtx.ID,
		Created:  utils.CurrentUnixTime(),
		Edited:   utils.CurrentUnixTime(),
		Rating:   0,
	}
	if newPost, status, err = repo.Create(&post, input.Categories); err != nil {
		response.Error(w, status, err)
		return
	}

	response.Success(w, fmt.Sprintf("post has been created"), newPost)
}

// UpdatePost modifies post record in database
//
// Returns Post model on success
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

	input.Validate(post)

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

// DeletePost removes records: post, it's reactions, it's categories,
// comments realted to that post and reactions to these comments from database
func DeletePost(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		postID int64
		post   *models.Post
		status int
		err    error
	)
	if postID, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if post, status, err = repo.FindByID(postID, -1); err != nil {
		response.Error(w, status, err)
		return
	}

	if !requestorIsEntityOwner(utils.GetUserFromCtx(r), post.AuthorID) {
		response.Error(w, http.StatusForbidden, errors.New("this post doesn't belong to you"))
		return
	}

	if err = crud.NewCommentRepoCRUD().DeleteGroup(postID); err != nil {
		logger.CheckErrAndLog("comments deletion", "", err)
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if err = crud.NewCategoryRepoCRUD().DeleteGroup(postID); err != nil {
		logger.CheckErrAndLog("categories deletion", "", err)
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if err = crud.NewPostRepoCRUD().DeleteAllReactions(postID); err != nil {
		logger.CheckErrAndLog("reactions deletion", "", err)
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if status, err = repo.Delete(postID); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, fmt.Sprintf("post has been deleted"), nil)
}
