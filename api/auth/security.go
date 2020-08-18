package auth

import (
	"fmt"
	"forum/config"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func generateCookie() *http.Cookie {
	return &http.Cookie{Name: config.SessionCookieName, Value: fmt.Sprint(uuid.NewV4()), Expires: time.Now().Add(config.CookieExpiration), HttpOnly: true}
}

//hash hashes the password
func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 8)
}

//verifyPassword checks if password is correct
func verifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
