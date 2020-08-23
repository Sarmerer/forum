package auth

import (
	"fmt"
	"forum/api/cache"
	"forum/config"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//generateCookie generates the cookie
func generateCookie() *http.Cookie {
	for {
		newUUID := fmt.Sprint(uuid.NewV4())
		if _, found := cache.Sessions.Get(newUUID); !found {
			return &http.Cookie{Name: config.SessionCookieName, Value: newUUID, Expires: time.Now().Add(config.CookieExpiration), HttpOnly: true}
		}
	}
}

//hash hashes the password
func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 8)
}

//verifyPassword verifyes the password
func verifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
