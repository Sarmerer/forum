package models

import (
	"forum/api/entities"
	"forum/database"
	"time"
)

const timeLayout = "2006-01-02 15:04:05"

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
		date, _ := time.Parse(timeLayout, created)
		user.Created = date
		date, _ = time.Parse(timeLayout, lastOnline)
		user.LastOnline = date
		users = append(users, user)
	}
	return users, nil
}

//Find returns a specific user from the database
func (*UserModel) Find(id int) (entities.User, error) {
	var user entities.User

	db, err := database.Connect()
	if err != nil {
		return user, err
	}
	rows, e := db.Query("SELECT * FROM users WHERE user_id = ?", id)
	if e != nil {
		return user, e
	}
	for rows.Next() {
		var created, lastOnline string
		rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Nickname, &created, &lastOnline, &user.SessionID, &user.Role)
		date, _ := time.Parse(timeLayout, created)
		user.Created = date
		date, _ = time.Parse(timeLayout, lastOnline)
		user.LastOnline = date
	}
	return user, nil
}

//Create adds a new user to the database
func (*UserModel) Create(user *entities.User) bool {
	db, err := database.Connect()
	statement, err := db.Prepare("INSERT INTO users (user_name, user_password, user_email, user_nickname, user_created, user_last_online, user_session_id, user_role) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return false
	}
	res, err := statement.Exec(user.Name, user.Password, user.Email, user.Nickname, time.Now().Format(timeLayout), time.Now().Format(timeLayout), user.SessionID, user.Role)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}

//Delete deletes user from the database
func (*UserModel) Delete(id int) bool {
	db, err := database.Connect()
	if err != nil {
		return false
	}
	res, err := db.Exec("DELETE FROM users WHERE user_id = ?", id)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}

//Update updates existing user in the database
func (*UserModel) Update(user *entities.User) bool {
	db, err := database.Connect()
	statement, err := db.Prepare("UPDATE users SET user_name = ?, user_password = ?, user_email = ?, user_nickname = ?, user_created = ?, user_last_online = ?, user_session_id = ?, user_role = ? WHERE user_id = ?")
	if err != nil {
		panic(err)
		return false
	}
	res, err := statement.Exec(user.Name, user.Password, user.Email, user.Nickname, user.Created.Format(timeLayout), user.LastOnline.Format(timeLayout), user.SessionID, user.Role, user.ID)
	if err != nil {
		panic(err)

		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		panic(err)

		return false
	}
	return rowsAffected > 0
}
