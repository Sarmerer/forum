package models

import (
	"errors"
	"regexp"
	"strings"
)

func (i *InputAllPosts) Validate() {
	m := make(map[string]bool)
	m["rating"] = true
	m["created"] = true
	m["total_participants"] = true
	m["comments_count"] = true
	if _, ok := m[i.OrderBy]; !ok {
		i.OrderBy = "rating"
	}

	var direction = strings.ToLower(i.Direction)

	if direction != "desc" && direction != "asc" {
		i.Direction = "asc"
	} else {
		i.Direction = direction
	}

	if i.CurrentPage < 1 {
		i.CurrentPage = 1
	}
	if i.PerPage < 1 {
		i.PerPage = 7
	}
}

func (i InputUserSignUp) Validate() error {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(i.Login) < 1 {
		return errors.New("login is too short")
	} else if len(i.Login) > 15 {
		return errors.New("login is too long")

	}
	if len(i.Password) < 5 {
		return errors.New("password is too short")
	}
	if !emailRegex.MatchString(i.Email) {
		return errors.New("invalid email")
	}
	return nil
}

func (i InputRate) Validate() error {
	if i.Reaction < -1 || i.Reaction > 1 {
		return errors.New("reaction should be -1, 0 or 1")
	}
	return nil
}

func (i InputPostCreateUpdate) Validate(post *Post) {
	if i.Title != "" {
		post.Title = i.Title
	}
	if i.Content != "" {
		post.Content = i.Content
	}

	post.IsImage = i.IsImage

}
