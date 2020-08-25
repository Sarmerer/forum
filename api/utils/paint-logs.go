package utils

import (
	"fmt"
	"net/http"
)

var (
	Green   = color("\033[97;42m%s\033[0m")
	White   = color("\033[90;47m%s\033[0m")
	Yellow  = color("\033[90;43m%s\033[0m")
	Red     = color("\033[97;41m%s\033[0m")
	Blue    = color("\033[97;44m%s\033[0m")
	Magenta = color("\033[97;45m%s\033[0m")
	Cyan    = color("\033[97;46m%s\033[0m")
)

func PaintMethod(method string) string {
	switch method {
	case http.MethodGet:
		return Blue(method + "    ")
	case http.MethodPost:
		return Green(method + "   ")
	case http.MethodPut:
		return Magenta(method + "    ")
	case http.MethodDelete:
		return Red(method + " ")
	case http.MethodOptions:
		return Yellow(method)
	default:
		return White(method)
	}
}

func PaintStatus(code int) string {
	switch {
	case code >= 0 && code < 300:
		return Green(code)
	case code >= 300 && code < 400:
		return Magenta(code)
	case code >= 400 && code < 500:
		return Yellow(code)
	case code >= 500:
		return Red(code)
	default:
		return White(code)
	}
}

func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
