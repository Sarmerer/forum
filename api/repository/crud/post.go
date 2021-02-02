package crud

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/utils"
)

//PostRepoCRUD helps performing CRUD operations
type PostRepoCRUD struct{}

//NewPostRepoCRUD creates an instance of PostModel
func NewPostRepoCRUD() PostRepoCRUD {
	return PostRepoCRUD{}
}

func (PostRepoCRUD) fetchCategories(post *models.Post) (status int, err error) {
	if post.Categories, status, err = NewCategoryRepoCRUD().FindByPostID(post.ID); err != nil {
		return status, err
	}
	return http.StatusOK, nil
}

func (PostRepoCRUD) fetchAuthor(post *models.Post) (status int, err error) {

	if post.Author, status, err = NewUserRepoCRUD().FindByID(post.AuthorID); err != nil {
		if status == http.StatusInternalServerError {
			return status, err
		}
		post.Author = DeletedUser
	}
	return http.StatusOK, nil
}

//FindAll returns all posts in the database
//TODO improve this MESS
func (PostRepoCRUD) FindAll(requestorID int64, input *models.InputAllPosts) (*models.Posts, int, error) {
	var (
		posts        *sql.Rows
		result       models.Posts
		offset       int    = (input.CurrentPage - 1) * input.PerPage
		recentAmount int    = 5
		orderBy      string = fmt.Sprintf("%s %s", input.OrderBy, input.Direction)
		status       int
		err          error
	)
	if posts, err = repository.DB.Query(
		fmt.Sprintf(
			`SELECT *,
			(
				SELECT TOTAL(reaction)
				FROM posts_reactions
				WHERE post_id_fkey = p._id
			) AS rating,
			IFNULL (
				(
					SELECT reaction
					FROM posts_reactions
					WHERE user_id_fkey = $1
						AND post_id_fkey = p._id
				),
				0
			) AS your_reaction,
			(
				SELECT count(_id)
				FROM comments
				WHERE post_id_fkey = p._id
				AND deleted = 0
			) AS comments_count,
			IFNULL (
				(
					SELECT COUNT(DISTINCT author_id_fkey)
					FROM comments
					WHERE post_id_fkey = p._id
					AND deleted = 0
				),
				0
			) AS total_participants
			FROM posts p
			ORDER BY %s
			LIMIT $2 OFFSET $3`,
			orderBy),
		requestorID, input.PerPage, offset,
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	for posts.Next() {
		var err error
		var p models.Post

		if err = posts.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.IsImage, &p.Created, &p.Edited, &p.EditReason,
			&p.Rating, &p.YourReaction, &p.CommentsCount,
			&p.ParticipantsCount,
		); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		if status, err = NewPostRepoCRUD().fetchAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = NewPostRepoCRUD().fetchCategories(&p); err != nil {
			return nil, status, err
		}
		result.Hot = append(result.Hot, p)
	}
	if err = repository.DB.QueryRow(`SELECT COUNT(_id) FROM posts`).Scan(&result.TotalRows); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if result.Recent, status, err = NewPostRepoCRUD().FindRecent(recentAmount); err != nil {
		return nil, status, err
	}
	return &result, http.StatusOK, nil
}

func (PostRepoCRUD) FindRecent(amount int) ([]models.Post, int, error) {
	var (
		posts  *sql.Rows
		result []models.Post
		status int
		err    error
	)
	if posts, err = repository.DB.Query(
		`SELECT *
		FROM posts
		ORDER BY created DESC
		LIMIT $1`,
		amount,
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	for posts.Next() {
		var p models.Post
		if err = posts.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.IsImage, &p.Created, &p.Edited, &p.EditReason,
		); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		if status, err = NewPostRepoCRUD().fetchAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = NewPostRepoCRUD().fetchCategories(&p); err != nil {
			return nil, status, err
		}
		if err = NewCommentRepoCRUD().CountForPost(&p); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		result = append(result, p)
	}
	return result, http.StatusOK, nil
}

//FindByID returns a specific post from the database
func (PostRepoCRUD) FindByID(postID int64, requestorID int64) (*models.Post, int, error) {
	var (
		p      models.Post
		status int
		err    error
	)
	if err = repository.DB.QueryRow(
		`SELECT *
		FROM posts
		WHERE _id = $1`,
		postID,
	).Scan(
		&p.ID, &p.AuthorID, &p.Title, &p.Content,
		&p.IsImage, &p.Created, &p.Edited, &p.EditReason,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, errors.New("post not found")
	}

	if status, err = NewPostRepoCRUD().fetchAuthor(&p); err != nil {
		return nil, status, err
	}
	if status, err = NewPostRepoCRUD().fetchCategories(&p); err != nil {
		return nil, status, err
	}
	if err = NewCommentRepoCRUD().CountForPost(&p); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if p.Rating, p.YourReaction, err = NewPostRepoCRUD().GetRating(p.ID, requestorID); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &p, http.StatusOK, nil
}

func (PostRepoCRUD) FindByAuthor(userID, requestorID int64) ([]models.Post, int, error) {
	var (
		rows   *sql.Rows
		posts  []models.Post
		status int
		err    error
	)
	if rows, err = repository.DB.Query(
		`SELECT *
		FROM posts
		WHERE author_id_fkey = $1`,
		userID,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, errors.New("user not found")
	}
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.IsImage, &p.Created, &p.Edited, &p.EditReason,
		)
		if status, err = NewPostRepoCRUD().fetchAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = NewPostRepoCRUD().fetchCategories(&p); err != nil {
			return nil, status, err
		}
		if err = NewCommentRepoCRUD().CountForPost(&p); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if p.Rating, p.YourReaction, err = NewPostRepoCRUD().GetRating(p.ID, requestorID); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		posts = append(posts, p)
	}
	return posts, http.StatusOK, nil
}

func (PostRepoCRUD) FindByCategories(categories []string, requestorID int64) ([]models.Post, int, error) {
	var (
		rows   *sql.Rows
		posts  []models.Post
		status int
		err    error
	)
	if rows, err = repository.DB.Query(
		fmt.Sprintf(`SELECT p.*
		FROM posts_categories_bridge AS pcb
		INNER JOIN posts as p ON p._id = pcb.post_id_fkey
		INNER JOIN categories AS c ON c._id = pcb.category_id_fkey
		WHERE c.name IN (%s)
		GROUP BY p._id
		HAVING COUNT(DISTINCT c._id) = $1`, fmt.Sprintf("\"%s\"", strings.Join(categories, "\", \""))),
		len(categories),
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content,
			&p.Created, &p.Edited, &p.EditReason,
		)

		if status, err = NewPostRepoCRUD().fetchAuthor(&p); err != nil {
			return nil, status, err
		}
		if status, err = NewPostRepoCRUD().fetchCategories(&p); err != nil {
			return nil, status, err
		}
		if err = NewCommentRepoCRUD().CountForPost(&p); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if p.Rating, p.YourReaction, err = NewPostRepoCRUD().GetRating(p.ID, requestorID); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		posts = append(posts, p)
	}
	return posts, http.StatusOK, nil
}

//Create adds a new post to the database
func (PostRepoCRUD) Create(post *models.Post, categories []string) (*models.Post, int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		now          int64 = utils.CurrentUnixTime()
		newPost      *models.Post
		status       int
		err          error
	)
	if result, err = repository.DB.Exec(
		`INSERT INTO posts (
			author_id_fkey,
			title,
			content,
			is_image,
			created,
			edited,
			edit_reason
		)
	VALUES (?, ?, ?, ?, ?, ?, ?)`,
		post.AuthorID, post.Title, post.Content,
		post.IsImage, now, now, post.EditReason,
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
			is_image = ?,
			created = ?,
			edited = ?,
			edit_reason = ?
		WHERE _id = ?`,
		post.AuthorID, post.Title, post.Content, post.IsImage, post.Created,
		utils.CurrentUnixTime(), post.EditReason, post.ID,
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
func (PostRepoCRUD) Delete(postID int64) (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`DELETE FROM posts
		WHERE _id = ?`, postID,
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
