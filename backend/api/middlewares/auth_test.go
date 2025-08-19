package middlewares_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/teams_register"
	"trxd/api/routes/users_register"
	"trxd/db/sqlc"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../", "auth")
}

var testAuthMiddlewares = []struct {
	method           string
	endpoint         string
	expectedStatuses []int
}{
	{
		method:           "POST",
		endpoint:         "/users/login",
		expectedStatuses: []int{http.StatusBadRequest, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden},
	},
	{
		method:           "GET",
		endpoint:         "/users/info",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusOK, http.StatusOK, http.StatusOK, http.StatusOK, http.StatusOK},
	},
	{
		method:           "GET",
		endpoint:         "/challenges",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusOK, http.StatusForbidden, http.StatusOK, http.StatusOK, http.StatusOK},
	},
	{
		method:           "POST",
		endpoint:         "/teams/register",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest},
	},
	{
		method:           "POST",
		endpoint:         "/categories/create",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusBadRequest, http.StatusBadRequest},
	},
	{
		method:           "PATCH",
		endpoint:         "/configs/update",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusBadRequest},
	},
}

func TestAuthMiddlewares(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	var err error
	users := [6]*sqlc.User{}
	users[1], err = users_register.RegisterUser(t.Context(), "spectator", "spectator@test.test", "testpass", sqlc.UserRoleSpectator)
	if err != nil {
		t.Fatalf("Failed to register spectator user: %v", err)
	}
	users[2], err = users_register.RegisterUser(t.Context(), "player", "player@test.test", "testpass", sqlc.UserRolePlayer)
	if err != nil {
		t.Fatalf("Failed to register player user: %v", err)
	}
	users[3], err = users_register.RegisterUser(t.Context(), "team_player", "team@test.test", "testpass", sqlc.UserRolePlayer)
	if err != nil {
		t.Fatalf("Failed to register player user: %v", err)
	}
	_, err = teams_register.RegisterTeam(t.Context(), "team1", "teampass", users[3].ID)
	if err != nil {
		t.Fatalf("Failed to register team: %v", err)
	}
	users[4], err = users_register.RegisterUser(t.Context(), "author", "author@test.test", "testpass", sqlc.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	users[5], err = users_register.RegisterUser(t.Context(), "admin", "admin@test.test", "testpass", sqlc.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}

	for _, test := range testAuthMiddlewares {
		for j, user := range users {
			session := test_utils.NewApiTestSession(t, app)
			if user != nil {
				session.Post("/users/login", JSON{"email": user.Email, "password": "testpass"}, http.StatusOK)
			}
			session.Request(test.method, test.endpoint, nil, test.expectedStatuses[j])
		}
	}
}
