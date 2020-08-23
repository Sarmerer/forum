package utils

import (
	"fmt"
	"net/http"
)

var (
	green   = color("\033[97;42m%s\033[0m")
	white   = color("\033[90;47m%s\033[0m")
	yellow  = color("\033[90;43m%s\033[0m")
	red     = color("\033[97;41m%s\033[0m")
	blue    = color("\033[97;44m%s\033[0m")
	magenta = color("\033[97;45m%s\033[0m")
	cyan    = color("\033[97;46m%s\033[0m")
)

func PaintMethod(method string) string {
	switch method {
	case http.MethodGet:
		return blue(method + "    ")
	case http.MethodPost:
		return green(method + "   ")
	case http.MethodPut:
		return magenta(method + "     ")
	case http.MethodDelete:
		return red(method + " ")
	case http.MethodOptions:
		return yellow(method)
	default:
		return white(method)
	}
}

func PaintStatus(code int) string {
	switch {
	case code >= 0 && code < 300:
		return green(code)
	case code >= 300 && code < 400:
		return magenta(code)
	case code >= 400 && code < 500:
		return yellow(code)
	case code >= 500:
		return red(code)
	default:
		return white(code)
	}
}

func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
