package middlewares_test

import (
	"net/http"
	"testing"
	"trxd/api"
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
		expectedStatuses: []int{http.StatusBadRequest, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden},
	},
	{
		method:           "GET",
		endpoint:         "/info",
		expectedStatuses: []int{http.StatusOK, http.StatusOK, http.StatusOK, http.StatusOK, http.StatusOK, http.StatusOK},
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
		endpoint:         "/categories",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusBadRequest, http.StatusBadRequest},
	},
	{
		method:           "PATCH",
		endpoint:         "/configs",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusBadRequest},
	},
}

func TestAuthMiddlewares(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer app.Shutdown()

	users := [6]*sqlc.User{}
	users[1] = test_utils.RegisterUser(t, "spectator", "spectator@test.test", "testpass", sqlc.UserRoleSpectator)
	users[2] = test_utils.RegisterUser(t, "player", "player@test.test", "testpass", sqlc.UserRolePlayer)
	users[3] = test_utils.RegisterUser(t, "team_player", "team@test.test", "testpass", sqlc.UserRolePlayer)
	test_utils.RegisterTeam(t, "team1", "teampass", users[3].ID)
	users[4] = test_utils.RegisterUser(t, "author", "author@test.test", "testpass", sqlc.UserRoleAuthor)
	users[5] = test_utils.RegisterUser(t, "admin", "admin@test.test", "testpass", sqlc.UserRoleAdmin)

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
