package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func ParseID(r *http.Request) (res int64, err error) {
	if res, err = strconv.ParseInt(r.URL.Query().Get("id"), 10, 64); err != nil {
		return 0, errors.New("invalid id")
	}
	return res, nil
}

func FormatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
