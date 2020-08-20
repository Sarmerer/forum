package main

import (
	"fmt"
	"forum/api/entities"
	"forum/api/models"
	"forum/config"
	"forum/database"
	"time"
)

func main() {
	tStart := time.Now()
	var forLoop, queryRow time.Duration
	for i := 0; i < 100; i++ {
		db, err := database.Connect()
		if err != nil {
			fmt.Println(err)
		}
		um, _ := models.NewUserModel(db)
		login := "lala11"
		var user entities.User
		rows, err := um.DB.Query("SELECT * FROM users WHERE user_name = ? OR user_email = ?", login, login)
		if err != nil {
			fmt.Println(err)
		}
		for rows.Next() {
			var created, lastOnline string
			rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
			date, _ := time.Parse(config.TimeLayout, created)
			user.Created = date
			date, _ = time.Parse(config.TimeLayout, lastOnline)
			user.LastOnline = date
		}
	}
	forLoop = time.Since(tStart)
	tStart = time.Now()
	for i := 0; i < 100; i++ {
		db, err := database.Connect()
		if err != nil {
			fmt.Println(err)
		}
		um, _ := models.NewUserModel(db)
		login := "lala11"
		var user entities.User
		var created, lastOnline string
		err = um.DB.QueryRow("SELECT * FROM users WHERE user_name = ? OR user_email = ?", login, login).Scan(
			&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role,
		)
		user.Created, _ = time.Parse(config.TimeLayout, created)
		user.LastOnline, _ = time.Parse(config.TimeLayout, lastOnline)
		if err != nil {
			fmt.Println(err)
		}
	}
	queryRow = time.Since(tStart)
	fmt.Printf("For loop took %v\t\tquery row took %v\n", forLoop, queryRow)
}
