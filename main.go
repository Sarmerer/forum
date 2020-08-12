package main

import (
	"fmt"
	"forum/api/models"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// api.Run()
	var um models.UserModel
	users, err := um.FindAll()
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}
