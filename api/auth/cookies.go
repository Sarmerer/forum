package auth

import (
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func generateCookie() *http.Cookie {
	return &http.Cookie{Name: "sessionID", Value: fmt.Sprint(uuid.NewV4()), Expires: time.Now().Add(14 * 24 * time.Hour), HttpOnly: true}
}
