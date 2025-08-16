package teams_get_test

import (
	"context"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/user_register"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "teams_get")
}

func TestTeamsGet(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	A, err := db.GetTeamByName(context.Background(), "A")
	if err != nil {
		t.Fatalf("Failed to get team A: %v", err)
	}
	if A == nil {
		t.Fatal("Team A not found")
	}
	B, err := db.GetTeamByName(context.Background(), "B")
	if err != nil {
		t.Fatalf("Failed to get team B: %v", err)
	}
	if B == nil {
		t.Fatal("Team B not found")
	}
	C, err := db.GetTeamByName(context.Background(), "C")
	if err != nil {
		t.Fatalf("Failed to get team C: %v", err)
	}
	if C == nil {
		t.Fatal("Team C not found")
	}

	expected := []map[string]interface{}{
		{
			"id":          A.ID,
			"name":        "A",
			"nationality": "",
			"score":       1498,
		},
		{
			"id":          B.ID,
			"name":        "B",
			"nationality": "",
			"score":       998,
		},
		{
			"id":          C.ID,
			"name":        "C",
			"nationality": "",
			"score":       0,
		},
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Get("/teams", nil, http.StatusOK)
	session.CheckResponse(expected)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"username": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/teams", nil, http.StatusOK)
	session.CheckResponse(expected)

	user, err := user_register.RegisterUser(context.Background(), "admin", "admin@admin.com", "adminpass", sqlc.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@admin.com", "password": "adminpass"}, http.StatusOK)
	session.Get("/teams", nil, http.StatusOK)
	session.CheckResponse(expected)
}
