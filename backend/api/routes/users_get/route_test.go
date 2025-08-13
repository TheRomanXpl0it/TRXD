package users_get_test

import (
	"context"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/user_register"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "user_get")
}

func TestUsersGet(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	expectedNoAuth := []JSON{
		{
			"email":       "",
			"image":       "",
			"name":        "a",
			"nationality": "",
			"role":        "",
			"score":       1498,
			"team_id":     nil,
		},
		{
			"email":       "",
			"image":       "",
			"name":        "b",
			"nationality": "",
			"role":        "",
			"score":       0,
			"team_id":     nil,
		},
		{
			"email":       "",
			"image":       "",
			"name":        "c",
			"nationality": "",
			"role":        "",
			"score":       998,
			"team_id":     nil,
		},
		{
			"email":       "",
			"image":       "",
			"name":        "d",
			"nationality": "",
			"role":        "",
			"score":       0,
			"team_id":     nil,
		},
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Get("/users", nil, http.StatusOK)
	body := session.Body()
	for _, user := range body.([]interface{}) {
		delete(user.(map[string]interface{}), "id")
	}
	err := utils.Compare(expectedNoAuth, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	expectedPlayer := []JSON{
		{
			"email":       "",
			"image":       "",
			"name":        "a",
			"nationality": "",
			"role":        "",
			"score":       1498,
			"team_id":     nil,
		},
		{
			"email":       "",
			"image":       "",
			"name":        "b",
			"nationality": "",
			"role":        "",
			"score":       0,
			"team_id":     nil,
		},
		{
			"email":       "",
			"image":       "",
			"name":        "c",
			"nationality": "",
			"role":        "",
			"score":       998,
			"team_id":     nil,
		},
		{
			"email":       "",
			"image":       "",
			"name":        "d",
			"nationality": "",
			"role":        "",
			"score":       0,
			"team_id":     nil,
		},
		{
			"email":       "",
			"image":       "",
			"name":        "test",
			"nationality": "",
			"role":        "",
			"score":       0,
			"team_id":     nil,
		},
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"username": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/users", nil, http.StatusOK)
	body = session.Body()
	for _, user := range body.([]interface{}) {
		delete(user.(map[string]interface{}), "id")
	}
	err = utils.Compare(expectedPlayer, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	expectedAdmin := []JSON{
		{
			"email":       "a@a",
			"image":       "",
			"name":        "a",
			"nationality": "",
			"role":        "Player",
			"score":       1498,
			"team_id":     nil,
		},
		{
			"email":       "b@b",
			"image":       "",
			"name":        "b",
			"nationality": "",
			"role":        "Player",
			"score":       0,
			"team_id":     nil,
		},
		{
			"email":       "c@c",
			"image":       "",
			"name":        "c",
			"nationality": "",
			"role":        "Player",
			"score":       998,
			"team_id":     nil,
		},
		{
			"email":       "d@d",
			"image":       "",
			"name":        "d",
			"nationality": "",
			"role":        "Player",
			"score":       0,
			"team_id":     nil,
		},
		{
			"email":       "e@e",
			"image":       "",
			"name":        "e",
			"nationality": "",
			"role":        "Admin",
			"score":       0,
			"team_id":     nil,
		},
		{
			"email":       "f@f",
			"image":       "",
			"name":        "f",
			"nationality": "",
			"role":        "Author",
			"score":       0,
			"team_id":     nil,
		},
		{
			"email":       "test@test.test",
			"image":       "",
			"name":        "test",
			"nationality": "",
			"role":        "Player",
			"score":       0,
			"team_id":     nil,
		},
		{
			"email":       "admin@test.com",
			"image":       "",
			"name":        "admin",
			"nationality": "",
			"role":        "Admin",
			"score":       0,
			"team_id":     nil,
		},
	}

	admin, err := user_register.RegisterUser(context.Background(), "admin", "admin@test.com", "testpass", db.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}
	if admin == nil {
		t.Fatal("Admin registration returned nil")
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get("/users", nil, http.StatusOK)
	body = session.Body()
	for _, user := range body.([]interface{}) {
		delete(user.(map[string]interface{}), "id")
	}
	err = utils.Compare(expectedAdmin, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}
}
