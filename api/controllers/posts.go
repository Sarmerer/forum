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

type postResponse struct {
	Post       *models.Post `json:"post"`
	Categories interface{}  `json:"categories"`
	Comments   interface{}  `json:"comments,omitempty"`
}

//TODO create SQL query that returns posts sorted by rating
func GetPosts(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		uid    int64
		posts  []models.Post
		result []postResponse
		err    error
	)
	uid = func() int64 {
		if r.Context().Value("uid") != nil {
			return r.Context().Value("uid").(int64)
		}
		return -1
	}()
	if posts, err = repo.FindAll(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	for i := 0; i < len(posts); i++ {
		p := postResponse{Post: &posts[i]}
		if p.Categories, err = GetCategoriesByPostID(posts[i].ID); err != nil {
			p.Categories = err
		}
		if p.Comments, err = CountComments(posts[i].ID); err != nil {
			p.Comments = err
		}
		if p.Post.Rating, p.Post.YourReaction, err = GetRating(p.Post.ID, uid); err != nil {
			fmt.Println(err)
		}
		result = append(result, p)
	}
	response.Success(w, nil, result)
}

func FindPost(w http.ResponseWriter, r *http.Request) {
	var (
		repo   repository.PostRepo = crud.NewPostRepoCRUD()
		input  models.InputPostFind
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
		result = append(result, p)
	}
	response.Success(w, nil, result)
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
		input  models.InputID
		status int
		err    error
	)
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	if _, status, err = repo.FindByID(input.ID); err != nil {
		response.Error(w, status, err)
		return
	}
	if err = DeleteCommentsGroup(input.ID); err != nil {
		fmt.Println(err)
	}
	if err = DeleteAllCategoriesForPost(input.ID); err != nil {
		fmt.Println(err)
	}
	if status, err = repo.Delete(input.ID); err != nil {
		response.Error(w, status, err)
		return
	}
	response.Success(w, fmt.Sprintf("post has been deleted"), nil)
}
