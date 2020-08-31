package crud

import (
	"database/sql"
)

//PostReaction struct contains info about a reaction to a post
type PostReaction struct {
	ID       int
	PostID   uint64
	UserID   uint64
	Reaction int
}

//PostReactionModel helps performing CRUD operations
type PostReactionModel struct {
	DB *sql.DB
}

//NewPostReactionModel creates an instance of PostReactionModel
func NewPostReactionModel(db *sql.DB) *PostReactionModel {
	return &PostReactionModel{db}
}

//FindAll returns all reactions in the database
func (um *PostReactionModel) FindAll() ([]PostReaction, error) {
	rows, e := um.DB.Query("SELECT * FROM posts_reaction")
	if e != nil {
		return nil, e
	}
	var reactions []PostReaction

	for rows.Next() {
		var reaction PostReaction
		rows.Scan(&reaction.ID, &reaction.PostID, &reaction.UserID, &reaction.Reaction)
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

//Find returns a specific reaction from the database
func (um *PostReactionModel) Find(id int) (PostReaction, error) {
	var reaction PostReaction
	rows, err := um.DB.Query("SELECT * FROM posts_reaction WHERE reaction_id = ?", id)
	if err != nil {
		return reaction, err
	}
	for rows.Next() {
		rows.Scan(&reaction.ID, &reaction.PostID, &reaction.UserID, &reaction.Reaction)
	}
	return reaction, nil
}

//Create adds a new reaction to the database
func (um *PostReactionModel) Create(reaction *PostReaction) (bool, string) {
	statement, err := um.DB.Prepare("INSERT INTO posts_reaction (post_id, user_id, reaction) VALUES (?, ?, ?)")
	if err != nil {
		return false, "Internal server error"
	}
	res, err := statement.Exec(reaction.PostID, reaction.UserID, reaction.Reaction)
	if err != nil {
		return false, err.Error()
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, "Internal server error"
	}
	if rowsAffected > 0 {
		return true, ""
	}
	return false, "Internal server error"
}

//Delete deletes reaction from the database
func (um *PostReactionModel) Delete(id int) bool {
	res, err := um.DB.Exec("DELETE FROM posts_reaction WHERE reaction_id = ?", id)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}

//Update updates existing reaction in the database
func (um *PostReactionModel) Update(reaction *PostReaction) bool {
	statement, err := um.DB.Prepare("UPDATE posts_reaction SET post_id = ?, user_id = ?, reaction = ? WHERE reaction_id = ?")
	if err != nil {
		return false
	}
	res, err := statement.Exec(reaction.PostID, reaction.UserID, reaction.Reaction, reaction.ID)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}
