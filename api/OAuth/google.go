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

type google struct {
	AccessTokenName string
}

type googleUser struct {
	Login  string `json:"name"`
	Avatar string `json:"picture"`
	Email  string `json:"email"`
}

var Google = google{AccessTokenName: "code"}

func (g google) Auth(query url.Values, sessionID string) (users []*models.User, status int, err error) {
	var (
		atr    *accessTokenResponse
		gUser  *googleUser
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		user   *models.User
		exists bool
	)
	if query.Get(g.AccessTokenName) == "" {
		return nil, http.StatusBadRequest, errors.New("access token not present")
	}
	if atr, err = g.getToken(query.Get(g.AccessTokenName)); err != nil {
		return
	}
	if gUser, err = g.getUser(atr); err != nil {
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

		if oldUser, status, err = repo.FindByLoginOrEmail([]string{gUser.Login, gUser.Email}); err != nil {
			return nil, status, err
		}

		if oldUser.OAuthProvider != ProviderGoogle {
			var merger *models.User = &models.User{}
			merger.Username = gUser.Login
			merger.Email = gUser.Email
			merger.Avatar = gUser.Avatar
			merger.OAuthProvider = ProviderGoogle
			return []*models.User{merger, oldUser}, http.StatusConflict, errors.New("conflict")
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

func (g google) getToken(code string) (atr *accessTokenResponse, err error) {
	var (
		req          *http.Request
		res          *http.Response
		uri          string       = "https://oauth2.googleapis.com/token"
		clientID     string       = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
		clientSecret string       = os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
		url          string       = fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=http://localhost:8081/auth?provider=google", uri, clientID, clientSecret, code)
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

func (g google) getUser(atr *accessTokenResponse) (user *googleUser, err error) {
	var (
		req    *http.Request
		res    *http.Response
		client *http.Client = &http.Client{}
	)
	if req, err = http.NewRequest("GET", fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?access_token=%s", atr.AccessToken), nil); err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", atr.TokenType, atr.AccessToken))
	req.Header.Add("Accept", "application/json")
	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, err
	}
	return
}
