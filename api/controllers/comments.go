package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"
)

// GetComments returns all comments for a single post as an array
func GetComments(w http.ResponseWriter, r *http.Request) {
	var (
		repo     repository.CommentRepo = crud.NewCommentRepoCRUD()
		userCtx  models.UserCtx         = utils.GetUserFromCtx(r)
		comments *models.Comments
		input    models.InputFindComments
		status   int
		err      error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if _, status, err = crud.NewPostRepoCRUD().FindByID(input.PostID, -1); err != nil {
		response.Error(w, status, err)
		return
	}
	if comments, err = repo.FindByPostID(&input, userCtx.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, nil, comments)
}

// FindComments allows to search comments by various parameters.
// Currently supported:
//
// - author - returns all comments from single user
func FindComments(w http.ResponseWriter, r *http.Request) {
	var (
		repo     repository.CommentRepo = crud.NewCommentRepoCRUD()
		userCtx  models.UserCtx         = utils.GetUserFromCtx(r)
		input    models.InputFindPost
		comments []models.Comment
		status   int
		err      error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	switch input.By {
	case "author":
		if comments, status, err = repo.FindByAuthor(input.AuthorID, userCtx.ID); err != nil {
			response.Error(w, status, err)
			return
		}
	default:
		response.Error(w, http.StatusBadRequest, errors.New("unknown search type"))
		return
	}
	response.Success(w, nil, comments)
}

// CreateComment adds new comment record to database.
//
// Returns Comment model on success
func CreateComment(w http.ResponseWriter, r *http.Request) {
	var (
		repo       repository.CommentRepo = crud.NewCommentRepoCRUD()
		userCtx    models.UserCtx         = utils.GetUserFromCtx(r)
		input      models.InputCommentCreateUpdate
		newComment *models.Comment
		status     int
		err        error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if _, status, err = crud.NewPostRepoCRUD().FindByID(input.PostID, -1); err != nil {
		response.Error(w, status, err)
		return
	}

	getLineage := func(parent models.Comment) string {
		if parent.ID == 0 {
			return "/"
		}
		if parent.Depth > 1 {
			return fmt.Sprintf("%s/%v", parent.Lineage, parent.ID)
		}
		return fmt.Sprintf("/%v", parent.ID)
	}

	comment := &models.Comment{
		Content:  input.Content,
		Created:  utils.CurrentUnixTime(),
		PostID:   input.PostID,
		ParentID: input.Parent.ID,
		Depth:    input.Parent.Depth + 1,
		Lineage:  getLineage(input.Parent),
		AuthorID: userCtx.ID,
	}
	if newComment, err = repo.Create(comment); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been added", newComment)
}

// UpdateComment modifies comment record in database with new data
//
// Returns Comment model on success
func UpdateComment(w http.ResponseWriter, r *http.Request) {
	var (
		repo           repository.CommentRepo = crud.NewCommentRepoCRUD()
		input          models.InputCommentCreateUpdate
		comment        *models.Comment
		updatedComment *models.Comment
		status         int
		err            error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if comment, status, err = repo.FindByID(input.ID); err != nil {
		response.Error(w, status, err)
		return
	}

	if !requestorIsEntityOwner(utils.GetUserFromCtx(r), comment.AuthorID) {
		response.Error(w, http.StatusForbidden, errors.New("this comment doesn't belong to you"))
		return
	}

	comment.Content = input.Content
	if updatedComment, err = repo.Update(comment); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been updated", updatedComment)
}

// DeleteComment removes comment and all it's reactions records from database
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	var (
		repo    repository.CommentRepo = crud.NewCommentRepoCRUD()
		cid     int64
		comment *models.Comment
		status  int
		err     error
	)
	if cid, err = utils.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if comment, status, err = repo.FindByID(cid); err != nil {
		response.Error(w, status, err)
		return
	}

	if !requestorIsEntityOwner(utils.GetUserFromCtx(r), comment.AuthorID) {
		response.Error(w, http.StatusForbidden, errors.New("this comment doesn't belong to you"))
		return
	}

	if err = repo.Delete(comment); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "reply has been deleted", nil)
}
