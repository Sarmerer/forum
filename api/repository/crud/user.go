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
	ID:       -1,
	Username: "deleted",
	Email:    "deleted",
	Avatar:   "https://avatars.dicebear.com/api/male/deleted",
	Alias:    "deleted",
	Role:     0,
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
		`SELECT COUNT(_id) as posts_count,
		(
			SELECT COUNT(_id)
			FROM comments
			WHERE author_id_fkey = $1
			AND deleted = 0
		) AS comments_count,
		(
			(
				SELECT TOTAL(reaction)
				FROM posts_reactions
				WHERE post_id_fkey IN (
						SELECT _id
						FROM posts
						WHERE author_id_fkey = $1
					)
			) + (
				SELECT TOTAL(reaction)
				FROM comments_reactions
				WHERE comment_id_fkey IN (
						SELECT _id
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
		rows.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.Avatar, &u.Alias,
			&u.Created, &u.LastActive, &u.SessionID, &u.Role, &u.Verified, &u.OAuthProvider,
		)
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
		`SELECT _id,
		 		username,
				email,
				avatar,
				alias,
				created,
				last_active,
				role,
				verified,
				oauth_provider
		FROM users
		WHERE _id = ?`, userID,
	).Scan(
		&u.ID, &u.Username, &u.Email, &u.Avatar, &u.Alias,
		&u.Created, &u.LastActive, &u.Role, &u.Verified, &u.OAuthProvider,
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
		now          int64 = utils.CurrentUnixTime()
		lastInsertID int64
		newUser      *models.User
		status       int
		err          error
	)
	name := repository.DB.QueryRow("SELECT username FROM users WHERE username = ?", user.Username).Scan(&user.Username)
	email := repository.DB.QueryRow("SELECT email FROM users WHERE email = ?", user.Email).Scan(&user.Email)
	if name == nil && email != nil {
		return nil, http.StatusConflict, errors.New("name is already taken")
	} else if email == nil && name != nil {
		return nil, http.StatusConflict, errors.New("email is already taken")
	} else if name == nil && email == nil {
		return nil, http.StatusConflict, errors.New("name and email are already taken")
	}

	if result, err = repository.DB.Exec(
		`INSERT INTO users (
			username,
			password,
			email,
			avatar,
			alias,
			created,
			last_active,
			session_id,
			role,
			verified,
			oauth_provider
		)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Username, user.Password, user.Email, user.Avatar, user.Alias, now, now,
		user.SessionID, user.Role, user.Verified, user.OAuthProvider,
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if newUser, status, err = NewUserRepoCRUD().FindByID(lastInsertID); err != nil {
		return nil, status, err
	}
	return newUser, http.StatusBadRequest, nil
}

//Update updates existing user in the database
//TODO decide what we will let users to update
func (UserRepoCRUD) Update(user *models.User) (*models.User, int, error) {
	var (
		updatedUser *models.User
		status      int
		err         error
	)
	if _, err = repository.DB.Exec(
		`UPDATE users
		SET username = $1,
			email = ?,
			avatar = ?,
			alias = $1,
			last_active = ?,
			oauth_provider = ?
		WHERE _id = ?`,
		user.Username, user.Email, user.Avatar,
		user.LastActive, user.OAuthProvider, user.ID,
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if updatedUser, status, err = NewUserRepoCRUD().FindByID(user.ID); err != nil {
		return nil, status, err
	}
	return updatedUser, http.StatusOK, nil
}

func (UserRepoCRUD) Verify(userID int64, newValue bool) error {
	var (
		err error
	)
	if _, err = repository.DB.Exec(
		`UPDATE users
		SET verified = $1
		WHERE _id = $2`,
		newValue, userID,
	); err != nil {
		return err
	}

	return nil
}

func (UserRepoCRUD) UpdateLastActivity(userID int64) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`UPDATE users
		SET last_active = ?
		WHERE _id = ?`,
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
		WHERE _id = ?`, userID,
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

//FindByLoginOrEmail finds a user by name or email in the database
func (UserRepoCRUD) FindByLoginOrEmail(username, email string) (*models.User, int, error) {
	var (
		u   models.User
		err error
	)
	if err = repository.DB.QueryRow(
		`SELECT _id,
				username,
	   			email,
	   			avatar,
				alias,
				created,
	   			last_active,
				role,
				verified,
				oauth_provider
		FROM users
		WHERE username = $1 OR email = $2 AND verified = 1`, username, email,
	).Scan(
		&u.ID, &u.Username, &u.Email, &u.Avatar, &u.Alias, &u.Created,
		&u.LastActive, &u.Role, &u.Verified, &u.OAuthProvider,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, errors.New("user not found")
	}
	return &u, http.StatusOK, nil
}

func (UserRepoCRUD) FindUnverifiedByEmail(email string) (*models.User, int, error) {
	var (
		u   models.User
		err error
	)
	if err = repository.DB.QueryRow(
		`SELECT _id,
				username,
	   			email,
	   			avatar,
				alias,
				created,
	   			last_active,
				role,
				verified,
				oauth_provider
		FROM users
		WHERE email = ? AND verified = 0`, email,
	).Scan(
		&u.ID, &u.Username, &u.Email, &u.Avatar, &u.Alias, &u.Created,
		&u.LastActive, &u.Role, &u.Verified, &u.OAuthProvider,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, errors.New("user not found")
	}
	return &u, http.StatusOK, nil
}

func (u UserRepoCRUD) Exists(username, email string) (exists bool, err error) {
	if err = repository.DB.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE username = $1 OR email = $2
			)`, username, email,
	).Scan(
		&exists,
	); err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
	}
	return exists, nil
}
