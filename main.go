package main

import (
	"fmt"
	"forum/api/models"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// api.Run()
	var um models.UserModel
	// users, err := um.FindAll()
	// if err != nil {
	// 	panic(err)
	// }
	// for _, user := range users {
	// 	fmt.Println(user)
	// }

	user, err := um.Find(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	// fmt.Println(um.Update(&user))

	// var user entities.User
	// user.Name = "blah"
	// user.Nickname = "blah"
	// user.Email = "blah"
	// user.Role = 0
	// user.Password = "qwerty12345"
	// user.SessionID = "somesession1"
	// fmt.Println(um.Create(&user))

	// fmt.Println(um.Delete(3))

}
