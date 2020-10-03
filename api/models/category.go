package models

//Category struct contains info about category
type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UseCount int64 `json:"use_count,omitempty"`
}
