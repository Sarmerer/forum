package controllers

import (
	"forum/api/config"
	"forum/api/models"
)

// requestorIsEntityOwner takes requestor's ID and role from the request context,
// and checks if requestor is entity owner, or he is an admin.
// If neither are true - API responds with HTTP status 403
func requestorIsEntityOwner(userCtx models.UserCtx, affectedUserID int64) bool {
	if userCtx.ID != affectedUserID && userCtx.Role < config.RoleAdmin {
		return false
	}
	return true
}
