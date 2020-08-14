package security

import (
	"forum/api/models"
	"forum/database"
)

//ValidateSession checks if the session id is valid
func ValidateSession(id string) (bool, error) {
	db, _ := database.Connect()
	um, _ := models.NewUserModel(db)
	return um.Validate(id)
}
