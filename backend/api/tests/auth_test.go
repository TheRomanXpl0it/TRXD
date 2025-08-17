package tests

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/user_register"
	"trxd/db/sqlc"
	"trxd/utils/test_utils"
)

var testAuthMiddlewares = []struct {
	method           string
	endpoint         string
	expectedStatuses []int
}{
	{
		method:           "POST",
		endpoint:         "/login",
		expectedStatuses: []int{http.StatusBadRequest, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden},
	},
	{
		method:           "GET",
		endpoint:         "/info",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusOK, http.StatusOK, http.StatusOK, http.StatusOK},
	},
	{
		method:           "POST",
		endpoint:         "/teams",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest},
	},
	{
		method:           "POST",
		endpoint:         "/category",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusForbidden, http.StatusBadRequest, http.StatusBadRequest},
	},
	{
		method:           "PATCH",
		endpoint:         "/config",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusBadRequest},
	},
}

func TestAuthMiddlewares(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	var err error
	users := [5]*sqlc.User{}
	users[1], err = user_register.RegisterUser(t.Context(), "spectator", "spectator@test.test", "testpass", sqlc.UserRoleSpectator)
	if err != nil {
		t.Fatalf("Failed to register spectator user: %v", err)
	}
	users[2], err = user_register.RegisterUser(t.Context(), "player", "player@test.test", "testpass", sqlc.UserRolePlayer)
	if err != nil {
		t.Fatalf("Failed to register player user: %v", err)
	}
	users[3], err = user_register.RegisterUser(t.Context(), "author", "author@test.test", "testpass", sqlc.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	users[4], err = user_register.RegisterUser(t.Context(), "admin", "admin@test.test", "testpass", sqlc.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}

	for _, test := range testAuthMiddlewares {
		for j, user := range users {
			session := test_utils.NewApiTestSession(t, app)
			if user != nil {
				session.Post("/login", JSON{"email": user.Email, "password": "testpass"}, http.StatusOK)
			}
			session.Request(test.method, test.endpoint, nil, test.expectedStatuses[j])
		}
	}
}
