package OAuth

import (
	"net/http"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/utils"
)

type oAuthUser struct {
	Login  string
	Avatar string
	Email  string
}

func createUser(from *oAuthUser, sessionID string) (*models.User, int, error) {
	var (
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		user   *models.User
		status int
		err    error
	)
	now := utils.CurrentUnixTime()
	newUser := &models.User{
		Username:      from.Login,
		Alias:         from.Login,
		Email:         from.Email,
		Avatar:        from.Avatar,
		Created:       now,
		LastActive:    now,
		Role:          config.RoleUser,
		SessionID:     sessionID,
		Verified:      true,
		OAuthProvider: ProviderGitHub,
	}
	if user, status, err = repo.Create(newUser); err != nil {
		return nil, status, err
	}
	return user, http.StatusOK, nil
}
func updateUser(from *oAuthUser, sessionID string) (*models.User, int, error) {
	var (
		repo   repository.UserRepo = crud.NewUserRepoCRUD()
		user   *models.User
		status int
		err    error
	)

	if user, status, err = repo.FindByLoginOrEmail([]string{from.Login, from.Email}); err != nil {
		return nil, status, err
	}

	user.Alias = from.Login
	user.Email = from.Email
	user.Avatar = from.Avatar
	user.LastActive = utils.CurrentUnixTime()
	user.OAuthProvider = ProviderGitHub

	if user, status, err = repo.Update(user); err != nil {
		return nil, status, err
	}
	if err = repo.UpdateSession(user.ID, sessionID); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return user, http.StatusOK, nil
}
