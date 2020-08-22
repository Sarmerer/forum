package utils

import (
	"fmt"
	"net/http"
)

var (
	get     = color("\033[1;44m%s\033[0m")
	post    = color("\033[1;42m%s\033[0m")
	put     = color("\033[1;45m%s\033[0m")
	delete  = color("\033[1;101m%s\033[0m")
	options = color("\033[1;103m%s\033[0m")
	def     = color("\033[1;45m%s\033[0m")
)

func Paint(method string) string {
	switch method {
	case http.MethodGet:
		return get(method)
	case http.MethodPost:
		return post(method)
	case http.MethodPut:
		return put(method)
	case http.MethodDelete:
		return delete(method)
	case http.MethodOptions:
		return options(method)
	default:
		return def(method)
	}
}

func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
