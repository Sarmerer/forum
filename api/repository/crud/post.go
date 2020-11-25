package crud

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
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
	LIMIT $2
	OFFSET $3`, input.OrderBy), userID, input.PerPage, offset,
	); err != nil {
		return nil, err
	}
	for posts.Next() {
		var err error
		var p models.Post
		posts.Scan(&p.ID, &p.AuthorID, &p.AuthorName, &p.Title, &p.Content, &p.Created,
			&p.Updated, &p.Rating, &p.YourReaction, &p.CommentsCount, &p.ParticipantsCount)
		if p.Categories, err = NewCategoryRepoCRUD().FindByPostID(p.ID); err != nil {
			p.Categories = append(p.Categories, models.Category{ID: 0, Name: err.Error()})
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
		posts.Scan(&p.ID, &p.AuthorID, &p.AuthorName, &p.Title, &p.Content, &p.Created, &p.Updated)
		result = append(result, p)
	}
	return result, nil
}

//FindByID returns a specific post from the database
func (PostRepoCRUD) FindByID(postID int64, userID int64) (*models.Post, int, error) {
	var (
		post models.Post
		err  error
	)
	// ? Query looks too bad
	if err = repository.DB.QueryRow(
		`SELECT id,
		author_id_fkey,
		author_name_fkey,
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
		c.*
		FROM posts p
		JOIN(
			SELECT COUNT(id) AS comments_count,
				COUNT(DISTINCT author_id_fkey) AS total_participants,
				IFNULL(last_comment_from_id, -1) AS last_comment_from_id,
				IFNULL(last_comment_from_name, "") AS last_comment_from_name,
				IFNULL(last_comment_date, "") AS last_comment_date
			FROM comments c
				JOIN(
					SELECT author_id_fkey AS last_comment_from_id,
						author_name_fkey AS last_comment_from_name,
						created AS last_comment_date
					FROM comments
					ORDER BY created
					LIMIT 1
				)
			WHERE post_id_fkey = $2
		) AS c
		WHERE p.id = $2`,
		userID, postID,
	).Scan(
		&post.ID, &post.AuthorID, &post.AuthorName, &post.Title, &post.Content,
		&post.Created, &post.Updated, &post.Rating, &post.YourReaction,
		&post.CommentsCount, &post.ParticipantsCount, &post.LastCommentFromID,
		&post.LastCommentFromName, &post.LastCommentDate,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, errors.New("post not found")
	}

	if post.Categories, err = NewCategoryRepoCRUD().FindByPostID(post.ID); err != nil {
		post.Categories = append(post.Categories, models.Category{ID: 0, Name: err.Error()})
	}
	return &post, http.StatusOK, nil
}

func (PostRepoCRUD) FindByAuthor(uid int64) ([]models.Post, error) {
	var (
		rows  *sql.Rows
		posts []models.Post
		err   error
	)
	if rows, err = repository.DB.Query(
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
		) AS yor_reaction,
		(
			SELECT count(id)
			FROM comments
			WHERE post_id_fkey = p.id
		) AS comments_count
	FROM posts p
	WHERE p.author_id_fkey = $1`,
		uid,
	); err != nil {
		return nil, err
	}
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, &p.AuthorName, &p.Title, &p.Content, &p.Created, &p.Updated, &p.Rating, &p.YourReaction, &p.CommentsCount)
		if p.Categories, err = NewCategoryRepoCRUD().FindByPostID(p.ID); err != nil {
			p.Categories = append(p.Categories, models.Category{ID: 0, Name: err.Error()})
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (PostRepoCRUD) FindByCategories(categories []string) ([]models.Post, error) {
	var (
		rows  *sql.Rows
		args  string
		posts []models.Post
		err   error
	)
	for _, category := range categories {
		if len(args) == 0 {
			args += fmt.Sprintf(`"%s"`, category)
		} else {
			args += fmt.Sprintf(`, "%s"`, category)
		}
	}
	if rows, err = repository.DB.Query(
		`SELECT p.id,
		p.author_id_fkey,
		p.author_name_fkey,
		p.title,
		p.content,
		p.created,
		p.updated
	FROM posts_categories_bridge AS pcb
		INNER JOIN posts as p ON p.id = pcb.post_id_fkey
		INNER JOIN categories AS c ON c.id = pcb.category_id_fkey
	WHERE c.name IN (`+args+`)
	GROUP BY p.id
	HAVING COUNT(DISTINCT c.id) = ?`,
		len(categories),
	); err != nil {
		return nil, err
	}
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.AuthorID, p.AuthorName, &p.Title, &p.Content, &p.Created, &p.Updated)
		posts = append(posts, p)
	}
	return posts, nil
}

//Create adds a new post to the database
func (PostRepoCRUD) Create(post *models.Post) (*models.Post, int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		pid          int64
		newPost      *models.Post
		status       int
		err          error
	)
	if result, err = repository.DB.Exec(
		`INSERT INTO posts (
			author_id_fkey,
			author_name_fkey,
			title,
			content,
			created,
			updated
		)
	VALUES (?, ?, ?, ?, ?, ?)`,
		post.AuthorID, post.AuthorName, post.Title, post.Content,
		time.Now().Format(config.TimeLayout), time.Now().Format(config.TimeLayout),
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if pid, err = result.LastInsertId(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if newPost, status, err = NewPostRepoCRUD().FindByID(pid, -1); err != nil {
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
func (PostRepoCRUD) Update(post *models.Post) error {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)
	if result, err = repository.DB.Exec(
		`UPDATE posts
		SET author_id_fkey = ?,
			author_name_fkey = ?,
			title = ?,
			content = ?,
			created = ?,
			updated = ?
		WHERE id = ?`,
		post.AuthorID, post.AuthorName, post.Title, post.Content, post.Created, time.Now().Format(config.TimeLayout), post.ID,
	); err != nil {
		return err
	}

	if rowsAffected, err = result.RowsAffected(); err == nil {
		return nil
	}
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("could not update the post")
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
