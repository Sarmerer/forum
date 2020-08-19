package auth_test

import (
	"forum/api/router"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
)

func TestSignIn(t *testing.T) {
	apitest.New().
		Handler(router.New()).
		Post("/signin").
		FormData("login", "lala").
		FormData("password", "lalal").
		Expect(t).
		Body(`{"status":"error", "code":401, "message":"user not authorized","data":null}`).
		Status(http.StatusUnauthorized).
		End()
}
