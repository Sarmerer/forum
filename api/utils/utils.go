package utils

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func ServeTemplate(path, layout string, w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseGlob(path+"*.html")).ExecuteTemplate(w, layout, http.StatusOK)
}

func ParseURLInt(path, prefix string) (int64, error) {
	if strings.HasPrefix(path, prefix) {
		return strconv.ParseInt(strings.TrimPrefix(path, prefix), 10, 64)
	}
	return 0, errors.New(fmt.Sprint(path, http.StatusNotFound))
}
