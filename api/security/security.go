package security

import (
	"forum/api/models"
	"forum/database"

	"golang.org/x/crypto/bcrypt"
)

//ValidateSession checks if the session id is valid
func ValidateSession(id string) (bool, error) {
	db, _ := database.Connect()
	um, _ := models.NewUserModel(db)
	return um.Validate(id)
}

//Hash hashes the password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 8)
}

//VerifyPassword checks if password is correct
func VerifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
