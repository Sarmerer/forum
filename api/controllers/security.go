package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sarmerer/forum/api/config"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func generateCookie(cookie *http.Cookie, err error) (string, string) {
	var newUUID string
	if err != nil {
		newUUID = fmt.Sprint(uuid.NewV4())
	} else {
		newUUID = cookie.Value
	}
	newCookie := &http.Cookie{
		Name:     config.SessionCookieName,
		Value:    newUUID,
		Expires:  time.Now().Add(config.SessionExpiration),
		Path:     "/",
		HttpOnly: true,
	}
	if config.Production {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	}
	return newCookie.String(), newUUID
}

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 8)
}

func verifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
