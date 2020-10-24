package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/sarmerer/forum/api/config"

	"github.com/sarmerer/forum/api/models"
)

func ParseID(r *http.Request) (res int64, err error) {
	if res, err = strconv.ParseInt(r.URL.Query().Get("id"), 10, 64); err != nil {
		return 0, errors.New("invalid id")
	}
	return res, nil
}

func GetUserFromCtx(r *http.Request) models.UserCtx {
	if r.Context().Value(config.UserCtxVarName) != nil {
		return r.Context().Value(config.UserCtxVarName).(models.UserCtx)
	}
	return models.UserCtx{ID: -1, Role: -1}
}

func FormatRequest(r *http.Request) string {
	var request []string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	return strings.Join(request, "\n")
}
