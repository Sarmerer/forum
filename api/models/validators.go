package models

import (
	"errors"
	"regexp"
)

func (i InputAllPosts) Validate() {
	m := make(map[string]bool)
	m["rating"] = true
	m["created"] = true
	m["total_participants"] = true
	if _, ok := m[i.OrderBy]; ok {
		if i.Ascending {
			i.OrderBy += "DESC"
		} else {
			i.OrderBy += "ASC"
		}
	} else {
		i.OrderBy = "rating DESC"
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
	if len(i.Login) < 2 {
		return errors.New("login is too short")
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