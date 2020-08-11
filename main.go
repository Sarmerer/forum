package main

import (
	"forum/api"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	api.Run()
}
