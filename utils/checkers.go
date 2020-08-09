package utils

import (
	"net/http"
	"strconv"
	"strings"
)

func ParseURL(path, prefix string) (interface{}, error) {
	if strings.HasPrefix(path, prefix) {
		return strconv.Atoi(strings.TrimPrefix(path, prefix))
	}
	return nil, &HTTPError{path, http.StatusNotFound}
}
