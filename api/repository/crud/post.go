package crud

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
)

//PostRepoCRUD helps performing CRUD operations
type PostRepoCRUD struct{}

//NewPostRepoCRUD creates an instance of PostModel
func NewPostRepoCRUD() PostRepoCRUD {
	return PostRepoCRUD{}
}

func fetchAuthorAndCategories(p *models.Post) error {
	var (
		status int
		err    error
	)
	if p.Categories, err = NewCategoryRepoCRUD().FindByPostID(p.ID); err != nil {
		p.Categories = append(p.Categories, models.Category{ID: 0, Name: err.Error()})
	}
	if p.Author, status, err = NewUserRepoCRUD().FindByID(p.AuthorID); err != nil {
		if status == http.StatusInternalServerError {
			return err
		}
		p.Author = DeletedUser
	}
	return nil
}

//FindAll returns all posts in the database
//TODO improve this MESS
func (PostRepoCRUD) FindAll(userID int64, input models.InputAllPosts) (*models.Posts, error) {
	var (
		posts        *sql.Rows
		result       models.Posts
		offset       int = (input.CurrentPage - 1) * input.PerPage
		recentAmount int = 5
		err          error
	)
	if posts, err = repository.DB.Query(
		fmt.Sprintf(
			`SELECT *,
			(
				SELECT TOTAL(reaction)
				FROM posts_reactions
				WHERE post_id_fkey = p.id
			) AS rating,
			IFNULL (
				(
					SELECT reaction
					FROM posts_reactions
					WHERE user_id_fkey = $1
						AND post_id_fkey = p.id
				),
				0
			) AS your_reaction,
			(
				SELECT count(id)
				FROM comments
				WHERE post_id_fkey = p.id
			) AS comments_count,
			IFNULL (
				(
					SELECT COUNT(DISTINCT author_id_fkey)
					FROM comments
					WHERE post_id_fkey = p.id
				),
				0
			) AS total_participants
			FROM posts p
			ORDER BY %s
			LIMIT $2 OFFSET $3`,
			input.OrderBy),
		userID, input.PerPage, offset,
	); err != nil {
		return nil, err
	}
	for posts.Next() {
		var err error
		var p models.Post
		posts.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content, &p.Created,
			&p.Updated, &p.Rating, &p.YourReaction, &p.CommentsCount, &p.ParticipantsCount)
		if err = fetchAuthorAndCategories(&p); err != nil {
			return nil, err
		}
		result.Hot = append(result.Hot, p)
	}
	if result.Recent, err = NewPostRepoCRUD().FindRecent(recentAmount); err != nil {
		return nil, err
	}
	if err = repository.DB.QueryRow(
		`SELECT COUNT(id) FROM posts`,
	).Scan(&result.TotalRows); err != nil {
		return nil, err
	}
	return &result, nil
}

func (PostRepoCRUD) FindRecent(amount int) ([]models.Post, error) {
	var (
		posts  *sql.Rows
		result []models.Post
		err    error
	)
	if posts, err = repository.DB.Query(
		`SELECT *
		FROM posts p
		ORDER BY created DESC
		LIMIT $1`,
		amount,
	); err != nil {
		return nil, err
	}
	for posts.Next() {
		var p models.Post
		posts.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content, &p.Created, &p.Updated)
		if err = fetchAuthorAndCategories(&p); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

//FindByID returns a specific post from the database
func (PostRepoCRUD) FindByID(postID int64, userID int64) (*models.Post, int, error) {
	var (
		p   models.Post
		err error
	)
	// ? Query looks too bad
	if err = repository.DB.QueryRow(
		`SELECT id,
		author_id_fkey,
		title,
		content,
		created,
		updated,
		(
			SELECT TOTAL(reaction)
			FROM posts_reactions
			WHERE post_id_fkey = p.id
		) AS rating,
		IFNULL (
			(
				SELECT reaction
				FROM posts_reactions
				WHERE user_id_fkey = $1
					AND post_id_fkey = p.id
			),
			0
		) AS your_reaction,
		(
			SELECT COUNT(id)
			FROM comments
			WHERE post_id_fkey = p.id
		) AS comments_count,
		IFNULL (
			(
				SELECT COUNT(DISTINCT author_id_fkey)
				FROM comments
				WHERE post_id_fkey = p.id
			),
			0
		) AS total_participants
		FROM posts p
		WHERE p.id = $2`,
		userID, postID,
	).Scan(
		&p.ID, &p.AuthorID, &p.Title, &p.Content,
		&p.Created, &p.Updated, &p.Rating, &p.YourReaction,
		&p.CommentsCount, &p.ParticipantsCount,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, errors.New("post not found")
	}

	if err = fetchAuthorAndCategories(&p); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &p, http.StatusOK, nil
}

func (PostRepoCRUD) FindByAuthor(userID, requestorID int64) ([]models.Post, int, error) {
	var (
		rows  *sql.Rows
		posts []models.Post
		err   error
	)
	if rows, err = repository.DB.Query(
		// FIXME fix yout_reaction subquery
		// it should take in requestorID instead of userID
		`SELECT *,
		(
			SELECT TOTAL(reaction)
			FROM posts_reactions
			WHERE post_id_fkey = p.id
		) AS rating,
		IFNULL (
			(
				SELECT reaction
				FROM posts_reactions
				WHERE user_id_fkey = $1
					AND post_id_fkey = p.id
			),
			0
		) AS your_reaction,
		(
			SELECT count(id)
			FROM comments
			WHERE post_id_fkey = p.id
		) AS comments_count,
		IFNULL (
			(
				SELECT COUNT(DISTINCT author_id_fkey)
				FROM comments
				WHERE post_id_fkey = p.id
			),
			0
		) AS total_participants
		FROM posts p
		WHERE p.author_id_fkey = $1`,
		userID, requestorID,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, errors.New("user not found")
	}
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content, &p.Created, &p.Updated, &p.Rating, &p.YourReaction, &p.CommentsCount, &p.ParticipantsCount)
		if err = fetchAuthorAndCategories(&p); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		posts = append(posts, p)
	}
	return posts, http.StatusOK, nil
}

func (PostRepoCRUD) FindByCategories(categories []string) ([]models.Post, error) {
	var (
		rows  *sql.Rows
		posts []models.Post
		err   error
	)
	if rows, err = repository.DB.Query(
		fmt.Sprintf(`SELECT p.*
	FROM posts_categories_bridge AS pcb
		INNER JOIN posts as p ON p.id = pcb.post_id_fkey
		INNER JOIN categories AS c ON c.id = pcb.category_id_fkey
	WHERE c.name IN (%s)
	GROUP BY p.id
	HAVING COUNT(DISTINCT c.id) = ?`, fmt.Sprintf("\"%s\"", strings.Join(categories, "\", \""))),
		len(categories),
	); err != nil {
		return nil, err
	}
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content, &p.Created, &p.Updated)
		if err = fetchAuthorAndCategories(&p); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

//Create adds a new post to the database
func (PostRepoCRUD) Create(post *models.Post, categories []string) (*models.Post, int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		newPost      *models.Post
		status       int
		err          error
	)
	if result, err = repository.DB.Exec(
		`INSERT INTO posts (
			author_id_fkey,
			title,
			content,
			created,
			updated
		)
	VALUES (?, ?, ?, ?, ?)`,
		post.AuthorID, post.Title, post.Content,
		time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout),
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if post.ID, err = result.LastInsertId(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(categories) > 0 {
		if err = NewCategoryRepoCRUD().Create(post.ID, categories); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}
	if newPost, status, err = NewPostRepoCRUD().FindByID(post.ID, -1); err != nil {
		return nil, status, err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return newPost, http.StatusOK, nil
	}
	return nil, http.StatusBadRequest, errors.New("could not create the post")
}

//Update updates existing post in the database
func (PostRepoCRUD) Update(post *models.Post, userCtx models.UserCtx) (*models.Post, int, error) {
	var (
		result       sql.Result
		updatedPost  *models.Post
		rowsAffected int64
		status       int
		err          error
	)
	if result, err = repository.DB.Exec(
		`UPDATE posts
		SET author_id_fkey = ?,
			title = ?,
			content = ?,
			created = ?,
			updated = ?
		WHERE id = ?`,
		post.AuthorID, post.Title, post.Content, post.Created, time.Now().Format(config.TimeLayout), post.ID,
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if updatedPost, status, err = NewPostRepoCRUD().FindByID(post.ID, userCtx.ID); err != nil {
		return nil, status, err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if rowsAffected > 0 {
		return updatedPost, http.StatusOK, nil
	}
	return nil, http.StatusBadRequest, errors.New("could not update the post")
}

//Delete deletes post from the database
func (PostRepoCRUD) Delete(pid int64) (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`DELETE FROM posts
		WHERE id = ?`, pid,
	); err != nil {
		if err != sql.ErrNoRows {
			return http.StatusInternalServerError, err
		}
		return http.StatusNotFound, errors.New("post not found")
	}

	if rowsAffected, err = result.RowsAffected(); err == nil {
		return http.StatusOK, nil
	}
	if rowsAffected > 0 {
		return http.StatusOK, nil
	}
	return http.StatusNotModified, errors.New("could not delete the post")
}
