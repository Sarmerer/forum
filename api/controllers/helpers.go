package controllers

import (
	"forum/api/config"
	"forum/api/utils"
	"net/http"
)

// requestorIsAccountOwner takes requestor's ID and role from the request context,
// and checks if requestor is an account owner, or he is an admin.
// If neither are true - API responds with HTTP status 403
func requestorIsAccountOwner(r *http.Request, affectedUserID int64) bool {
	var userCtx = utils.GetUserFromCtx(r)
	if userCtx.ID != affectedUserID && userCtx.Role < config.RoleAdmin {
		return false
	}
	return true
}

func requestorIsEntityOwner() {

}
