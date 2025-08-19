package teams_scoreboard_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/users_register"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "teams_scoreboard")
}

func TestTeamsScoreboard(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	A, err := db.GetTeamByName(t.Context(), "A")
	if err != nil {
		t.Fatalf("Failed to get team A: %v", err)
	}
	if A == nil {
		t.Fatal("Team A not found")
	}
	B, err := db.GetTeamByName(t.Context(), "B")
	if err != nil {
		t.Fatalf("Failed to get team B: %v", err)
	}
	if B == nil {
		t.Fatal("Team B not found")
	}
	C, err := db.GetTeamByName(t.Context(), "C")
	if err != nil {
		t.Fatalf("Failed to get team C: %v", err)
	}
	if C == nil {
		t.Fatal("Team C not found")
	}

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

	user, err := users_register.RegisterUser(t.Context(), "admin", "admin@admin.com", "adminpass", sqlc.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "admin@admin.com", "password": "adminpass"}, http.StatusOK)
	session.Get("/teams/scoreboard", nil, http.StatusOK)
	session.CheckResponse(expected)
}
