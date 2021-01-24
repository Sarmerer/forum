package OAuth

import (
	"net/url"

	"github.com/sarmerer/forum/api/models"
)

// Available providers
var (
	ProviderGitHub string = "github"
	ProviderGoogle string = "google"
)

type Provider interface {
	Auth(query url.Values, sessionID string) (user []*models.User, status int, err error)
}

var Providers = map[string]Provider{
	"github": GitHub,
	"google": Google,
}
