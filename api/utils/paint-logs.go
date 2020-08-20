package utils

import "fmt"

var (
	GET   = color("\033[1;44;30m%s\033[0m")
	POST     = color("\033[1;42;30m%s\033[0m")
	PUT   = color("\033[1;45;30m%s\033[0m")
	DELETE  = color("\033[1;101;30m%s\033[0m")
	OPTIONS  = color("\033[1;103;30m%s\033[0m")
	Default = color("\033[1;45;30m%s\033[0m")
)

func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
