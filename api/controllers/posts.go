package controllers

import (
	"encoding/json"
	"fmt"
	"forum/api/helpers"
	"forum/api/models"
	"forum/api/repository"
	"forum/api/response"
	"forum/config"
	"net/http"
	"time"
)

//TODO improve response time
func GetPosts(w http.ResponseWriter, r *http.Request) {
	type postTpl struct {
		Post       models.Post
		Categories interface{}
		Replies    interface{}
	}
	var (
		pm    repository.PostRepo
		posts []models.Post
		res   []postTpl
		err   error
	)
	if pm, err = helpers.PreparePostRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if posts, err = pm.FindAll(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	for _, post := range posts {
		p := postTpl{Post: post}
		if p.Replies, err = CountReplies(post.ID); err != nil {
			p.Replies = err
		}
		if p.Categories, err = GetCategories(post.ID); err != nil {
			p.Categories = err
		}
		res = append(res, p)
	}
	response.Success(w, nil, res)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	var (
		pid        uint64
		pm         repository.PostRepo
		post       *models.Post
		replies    []models.PostReply
		categories []models.Category
		status     int
		err        error
	)
	if pid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if pm, err = helpers.PreparePostRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if post, status, err = pm.FindByID(pid); err != nil {
		response.Error(w, status, err)
		return
	}

	if replies, err = GetReplies(pid); err != nil {
		replies = nil
	}
	if categories, err = GetCategories(pid); err != nil {
		categories = nil
	}
	res := struct {
		Post       interface{} `json:"post"`
		Categories interface{} `json:"categories"`
		Replies    interface{} `json:"replies"`
	}{post, categories, replies}
	response.Success(w, nil, res)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var (
		author            uint64
		pm                repository.PostRepo
		post              models.Post
		pid               int64
		cm                repository.CategoryRepo
		categoriesCreated int
		err               error
	)
	author = r.Context().Value("uid").(uint64)
	input := struct {
		Description string   `json:"description"`
		Content     string   `json:"content"`
		Categories  []string `json:"categories"`
	}{}

	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if pm, err = helpers.PreparePostRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	post = models.Post{
		Title:   input.Description,
		Content: input.Content,
		Author:  author,
		Created: time.Now().Format(config.TimeLayout),
		Updated: time.Now().Format(config.TimeLayout),
		Rating:  0,
	}
	if pid, err = pm.Create(&post); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if len(input.Categories) > 0 {
		if cm, err = helpers.PrepareCategoriesRepo(); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		if categoriesCreated, err = cm.Create(pid, input.Categories...); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
	}

	response.Success(w, fmt.Sprintf("post has been created. Catgories created: %d", categoriesCreated), nil)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var (
		name        string
		content     string
		pid         uint64
		pm          repository.PostRepo
		updatedPost *models.Post
		status      int
		err         error
	)
	name = r.FormValue("description")
	content = r.FormValue("content")
	if pid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if pm, err = helpers.PreparePostRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if updatedPost, status, err = pm.FindByID(pid); err != nil {
		response.Error(w, status, err)
		return
	}

	updatedPost.Updated = time.Now().Format(config.TimeLayout)
	if name != "" {
		updatedPost.Title = name
	}
	if content != "" {
		updatedPost.Content = content
	}

	if err = pm.Update(updatedPost); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, fmt.Sprint("post has been updated"), nil)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	var (
		pid    uint64
		pm     repository.PostRepo
		status int
		err    error
	)
	if pid, err = helpers.ParseID(r); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if pm, err = helpers.PreparePostRepo(); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	if _, status, err = pm.FindByID(pid); err != nil {
		response.Error(w, status, err)
		return
	}
	if status, err = pm.Delete(pid); err != nil {
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
