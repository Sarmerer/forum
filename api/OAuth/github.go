package OAuth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/utils"
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

func (gh gitHub) Auth(query url.Values, sessionID string) (users []*models.User, status int, err error) {
	var (
		atr    *accessTokenResponse
		gUser  *gitHubUser
		user   *models.User
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		exists bool
	)

	if query.Get(gh.AccessTokenName) == "" {
		return nil, http.StatusBadRequest, errors.New("access token not present")
	}
	if atr, err = getToken(query.Get(gh.AccessTokenName)); err != nil {
		return
	}
	if gUser, err = getUser(atr.AccessToken); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if gUser.Login == "" || gUser.Email == "" {
		return nil, http.StatusBadRequest, errors.New("invalid or expired token")
	}
	if exists, err = repo.Exists([]string{gUser.Login, gUser.Email}); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if exists {
		var oldUser *models.User

		if oldUser, status, err = repo.FindByLoginOrEmail(gUser.Login); err != nil {
			return nil, status, err
		}

		if oldUser.OAuthProvider != ProviderGitHub {
			var mergant *models.User = &models.User{}
			mergant.Username = gUser.Login
			mergant.Email = gUser.Email
			mergant.Avatar = gUser.Avatar
			mergant.OAuthProvider = ProviderGitHub
			return []*models.User{mergant, oldUser}, http.StatusConflict, errors.New("conflict")
		}

		oldUser.Username = gUser.Login
		oldUser.Email = gUser.Email
		oldUser.Avatar = gUser.Avatar
		oldUser.LastActive = utils.CurrentUnixTime()
		oldUser.OAuthProvider = ProviderGitHub

		if user, status, err = repo.Update(oldUser); err != nil {
			return nil, status, err
		}
		if err = repo.UpdateSession(oldUser.ID, sessionID); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		users = append(users, user)
	} else {
		now := utils.CurrentUnixTime()
		newUser := &models.User{
			Username:      gUser.Login,
			Alias:         gUser.Login,
			Email:         gUser.Email,
			Avatar:        gUser.Avatar,
			Created:       now,
			LastActive:    now,
			Role:          config.RoleUser,
			SessionID:     sessionID,
			OAuthProvider: ProviderGitHub,
		}
		if user, status, err = repo.Create(newUser); err != nil {
			return nil, status, err
		}
		users = append(users, user)
	}
	return users, http.StatusOK, nil
}

func getToken(code string) (atr *accessTokenResponse, err error) {
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

func getUser(token string) (user *gitHubUser, err error) {
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
