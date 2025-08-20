package teams_scoreboard_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	A := test_utils.GetTeamByName(t, "A")
	B := test_utils.GetTeamByName(t, "B")
	C := test_utils.GetTeamByName(t, "C")

	expected := []JSON{
		{
			"badges": []JSON{
				{
					"description": "Completed all cat-1 challenges",
					"name":        "cat-1",
				},
			},
			"country": "",
			"id":      A.ID,
			"name":    "A",
			"score":   1498,
		},
		{
			"badges": []JSON{
				{
					"description": "Completed all cat-2 challenges",
					"name":        "cat-2",
				},
			},
			"country": "",
			"id":      B.ID,
			"name":    "B",
			"score":   998,
		},
		{
			"country": "",
			"id":      C.ID,
			"name":    "C",
			"score":   0,
		},
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Get("/teams/scoreboard", nil, http.StatusOK)
	session.CheckResponse(expected)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/register", JSON{"username": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/teams/scoreboard", nil, http.StatusOK)
	session.CheckResponse(expected)

	test_utils.RegisterUser(t, "admin", "admin@test.test", "adminpass", sqlc.UserRoleAdmin)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "admin@test.test", "password": "adminpass"}, http.StatusOK)
	session.Get("/teams/scoreboard", nil, http.StatusOK)
	session.CheckResponse(expected)
}
