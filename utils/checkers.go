package utils

import (
	"net/http"
	"strconv"
	"strings"
)

func ParseURL(path, prefix string) (interface{}, error) {
	if strings.HasPrefix(path, prefix) {
		return strconv.ParseInt(strings.TrimPrefix(path, prefix), 10, 64)
	}
	return nil, &HTTPError{path, http.StatusNotFound}
}
