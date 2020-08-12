package models

import (
	"forum/api/entities"
	"forum/database"
	"time"
)

//UserModel helps performing CRUD operations
type UserModel struct {
}

//FindAll returns all users in the database
func (*UserModel) FindAll() ([]entities.User, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	rows, e := db.Query("SELECT * FROM users")
	if e != nil {
		return nil, e
	}
	var users []entities.User

	for rows.Next() {
		var user entities.User
		var created, lastOnline string
		rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
		date, _ := time.Parse("2006-01-02 15:04:05", created)
		user.Created = date
		date, _ = time.Parse("2006-01-02 15:04:05", lastOnline)
		user.LastOnline = date
		users = append(users, user)
	}
	return users, nil
}
