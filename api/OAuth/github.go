package OAuth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
)

type gitHub struct {
	AccessTokenName string
}

type gitHubUser struct {
	Login  string `json:"login"`
	Avatar string `json:"avatar_url"`
	Email  string `json:"email"`
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

var GitHub = gitHub{AccessTokenName: "code"}

func (gh gitHub) Auth(query url.Values, sessionID string) (user *models.User, status int, err error) {
	var (
		atr    *accessTokenResponse
		ghUser *gitHubUser
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		exists bool
	)

	if query.Get(gh.AccessTokenName) == "" {
		return nil, http.StatusBadRequest, errors.New("access token not present")
	}
	if atr, err = gh.getToken(query.Get(gh.AccessTokenName)); err != nil {
		return
	}
	if ghUser, err = gh.getUser(atr.AccessToken); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if ghUser.Login == "" || ghUser.Email == "" {
		return nil, http.StatusBadRequest, errors.New("invalid or expired token")
	}
	if exists, err = repo.Exists([]string{ghUser.Login, ghUser.Email}); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if exists {
		if user, status, err = updateUser((*oAuthUser)(ghUser), sessionID); err != nil {
			return nil, status, err
		}
	} else {
		if user, status, err = createUser((*oAuthUser)(ghUser), sessionID); err != nil {
			return nil, status, err
		}
	}
	return user, http.StatusOK, nil
}

func (gh gitHub) getToken(code string) (atr *accessTokenResponse, err error) {
	var (
		req          *http.Request
		res          *http.Response
		uri          string       = "https://github.com/login/oauth/access_token"
		clientID     string       = os.Getenv("GITHUB_OAUTH_CLIENT_ID")
		clientSecret string       = os.Getenv("GITHUB_OAUTH_CLIENT_SECRET")
		url          string       = fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s", uri, clientID, clientSecret, code)
		client       *http.Client = &http.Client{}
	)
	if req, err = http.NewRequest("POST", url, nil); err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	if res, err = client.Do(req); err != nil {
		return nil, err
	}
	if err = json.NewDecoder(res.Body).Decode(&atr); err != nil {
		return nil, err
	}
	return
}

func (gh gitHub) getUser(token string) (user *gitHubUser, err error) {
	var (
		req    *http.Request
		res    *http.Response
		client *http.Client = &http.Client{}
	)
	if req, err = http.NewRequest("GET", "https://api.github.com/user", nil); err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	if res, err = client.Do(req); err != nil {
		return nil, err
	}
	if err = json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, err
	}
	return
}
