package auth

import (
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

var activeSessions map[uuid.UUID]uuid.UUID = make(map[uuid.UUID]uuid.UUID)

func generateCookie() *http.Cookie {
	var u uuid.UUID
	for {
		u = uuid.NewV4()
		if _, found := activeSessions[u]; !found {
			activeSessions[u] = u
			break
		}
	}
	return &http.Cookie{Name: "sessionID", Value: fmt.Sprint(u), Expires: time.Now().Add(14 * 24 * time.Hour), HttpOnly: true}
}
