package middleware

import (
	"net/http"

	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	"github.com/sarmerer/forum/api/response"
	"github.com/sarmerer/forum/api/utils"
)

// UpdateUserActivity middleware is attached to routes that
// indicate that user is avtive, e.g. /post/new post/update...
// This function updates last activity date of a user in database
func UpdateUserActivity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			repo    repository.UserRepo = crud.NewUserRepoCRUD()
			userCtx models.UserCtx      = utils.GetUserFromCtx(r)
			err     error
		)
		if err = repo.UpdateLastActivity(userCtx.ID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		next(w, r)
	}
}
