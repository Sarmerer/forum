package controllers

import (
	"fmt"
	"forum/api/config"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//generateUUID generates the cookie
func generateCookie(cookie *http.Cookie, err error) *http.Cookie {
	var newUUID string
	if err != nil {
		newUUID = fmt.Sprint(uuid.NewV4())
	} else {
		newUUID = cookie.Value
	}
	return &http.Cookie{
		Name:     config.SessionCookieName,
		Value:    newUUID,
		Expires:  time.Now().Add(config.SessionExpiration),
		Path:     "/",
		HttpOnly: true}
}

//hash hashes the password
func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 8)
}

//verifyPassword verifyes the password
func verifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
