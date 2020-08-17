package auth

import (
	"fmt"
	"forum/config"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func generateCookie() *http.Cookie {
	return &http.Cookie{Name: config.SessionCookieName, Value: fmt.Sprint(uuid.NewV4()), Expires: time.Now().Add(config.CookieExpiration), HttpOnly: true}
}
