package session

import (
	"forum/api/models"
	"forum/database"
)

func New() {
}

func Destroy() {

}

//Validate checks if the session id exists in the database,
//which will mean that user is authorized
func Validate(sessionID string) (bool, error) {
	db, dbErr := database.Connect()
	defer db.Close()
	um, umErr := models.NewUserModel(db)
	if dbErr != nil {
		return false, dbErr
	}
	if umErr != nil {
		return false, umErr
	}
	return um.Validate(sessionID)
}
