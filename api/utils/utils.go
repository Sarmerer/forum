package utils

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func ServeTemplate(path, layout string, w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseGlob(path+"*.html")).ExecuteTemplate(w, layout, http.StatusOK)
}

func ParseURLInt(path, prefix string) (int64, error) {
	u, err := url.Parse(path)
	if err != nil {
		return 0, errors.New(fmt.Sprint(path, http.StatusInternalServerError))
	}
	u.RawQuery = ""
	return strconv.ParseInt(strings.TrimPrefix(u.String(), prefix), 10, 64)

}
