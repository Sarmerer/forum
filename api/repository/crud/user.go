package crud

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/utils"
)

//UserRepoCRUD helps performing CRUD operations
type UserRepoCRUD struct{}

var DeletedUser = &models.User{
	ID:          -1,
	Login:       "deleted",
	Email:       "deleted",
	Avatar:      "https://avatars.dicebear.com/api/male/deleted",
	DisplayName: "deleted",
	Role:        0,
}

//NewUserRepoCRUD creates an instance of UserModel
func NewUserRepoCRUD() UserRepoCRUD {
	return UserRepoCRUD{}
}

func (UserRepoCRUD) fetchUserStats(user *models.User) error {
	var (
		err error
	)
	if err = repository.DB.QueryRow(
		`SELECT COUNT(id) as posts_count,
		(
			SELECT COUNT(id)
			FROM comments
			WHERE author_id_fkey = $1
		) AS comments_count,
		(
			(
				SELECT TOTAL(reaction)
				FROM posts_reactions
				WHERE post_id_fkey IN (
						SELECT id
						FROM posts
						WHERE author_id_fkey = $1
					)
			) + (
				SELECT TOTAL(reaction)
				FROM comments_reactions
				WHERE comment_id_fkey IN (
						SELECT id
						FROM comments
						WHERE author_id_fkey = $1
					)
			)
		) AS rating
		FROM posts
		WHERE author_id_fkey = $1`,
		user.ID,
	).Scan(&user.Posts, &user.Comments, &user.Rating); err != nil {
		if err == sql.ErrNoRows {
			return err
		}
	}
	return nil
}

//FindAll returns all users in the database
func (UserRepoCRUD) FindAll() ([]models.User, error) {
	var (
		rows  *sql.Rows
		users []models.User
		err   error
	)
	if rows, err = repository.DB.Query(
		`SELECT *
		FROM users`,
	); err != nil {
		return nil, err
	}
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Login, &u.Password, &u.Email, &u.Avatar, &u.DisplayName, &u.Created, &u.LastActive, &u.SessionID, &u.Role)
		users = append(users, u)
	}
	return users, nil
}

//FindByID returns a specific user from the database
func (UserRepoCRUD) FindByID(userID int64) (*models.User, int, error) {
	var (
		u   models.User
		err error
	)
	if err = repository.DB.QueryRow(
		`SELECT id,
		 		login,
				email,
				avatar,
				display_name,
				created,
				last_online,
				role
		FROM users
		WHERE id = ?`, userID,
	).Scan(
		&u.ID, &u.Login, &u.Email, &u.Avatar, &u.DisplayName, &u.Created, &u.LastActive, &u.Role,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, errors.New("user not found")
	}
	if err = NewUserRepoCRUD().fetchUserStats(&u); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &u, http.StatusOK, nil
}

//Create adds a new user to the database
func (UserRepoCRUD) Create(user *models.User) (*models.User, int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		now          int64 = utils.CurrentUnixTime()
		lastInsertID int64
		newUser      *models.User
		status       int
		err          error
	)
	name := repository.DB.QueryRow("SELECT login FROM users WHERE login = ?", user.Login).Scan(&user.Login)
	email := repository.DB.QueryRow("SELECT email FROM users WHERE email = ?", user.Email).Scan(&user.Email)
	if name == nil && email != nil {
		return nil, http.StatusConflict, errors.New("name is not unique")
	} else if email == nil && name != nil {
		return nil, http.StatusConflict, errors.New("email is not unique")
	} else if name == nil && email == nil {
		return nil, http.StatusConflict, errors.New("name and email are not unique")
	}

	if result, err = repository.DB.Exec(
		`INSERT INTO users (
			login,
			password,
			email,
			avatar,
			display_name,
			created,
			last_online,
			session_id,
			role
		)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Login, user.Password, user.Email, user.Avatar, user.DisplayName, now, now, user.SessionID, user.Role,
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if newUser, status, err = NewUserRepoCRUD().FindByID(lastInsertID); err != nil {
		return nil, status, err
	}
	if rowsAffected > 0 {
		return newUser, http.StatusOK, nil
	}
	return nil, http.StatusBadRequest, errors.New("could not create the user")
}

//Update updates existing user in the database
//TODO decide what we will let users to update
func (UserRepoCRUD) Update(user *models.User) (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`UPDATE users
		SET display_name = ?
		WHERE id = ?`,
		user.DisplayName, user.ID,
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not update the user")
}

func (UserRepoCRUD) UpdateLastActivity(userID int64) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`UPDATE users
		SET last_online = ?
		WHERE id = ?`,
		utils.CurrentUnixTime(), userID,
	); err != nil {
		return err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("could not update user activity")
}

//Delete deletes user from the database
func (UserRepoCRUD) Delete(userID int64) (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`DELETE FROM users
		WHERE id = ?`, userID,
	); err != nil {
		if err != sql.ErrNoRows {
			return http.StatusInternalServerError, err
		}
		return http.StatusNotFound, errors.New("user not found")
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not delete the user")
}

//FindByNameOrEmail finds a user by name or email in the database
func (UserRepoCRUD) FindByLoginOrEmail(login string) (*models.User, int, error) {
	var (
		u   models.User
		err error
	)
	if err = repository.DB.QueryRow(
		`SELECT id,
				login,
	   			email,
	   			avatar,
				   display_name,
				   created,
	   			last_online,
	   			role
		FROM users
		WHERE login = $1
			OR email = $1`, login,
	).Scan(
		&u.ID, &u.Login, &u.Email, &u.Avatar, &u.DisplayName, &u.Created, &u.LastActive, &u.Role,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, errors.New("user not found")
	}
	return &u, http.StatusOK, nil
}
