package auth

import (
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func generateCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "sessionID",
		Value:    fmt.Sprint(uuid.NewV4()),
		MaxAge:   10 * 60,
		HttpOnly: true,
	}
}
