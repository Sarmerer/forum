package OAuth

import (
	"net/http"
	"net/url"

	"github.com/sarmerer/forum/api/models"
)

type google struct {
	AccessTokenName string
}

var Google = google{}

func (g google) Auth(query url.Values, sessionID string) (user *models.User, status int, err error) {
	return nil, http.StatusOK, nil
}
