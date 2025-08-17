package users_get_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/user_register"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "users_get")
}

func TestUsersGet(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	expectedNoAuth := []JSON{
		{
			"country": "",
			"email":   "",
			"image":   "",
			"name":    "a",
			"role":    "",
			"score":   1498,
		},
		{
			"country": "",
			"email":   "",
			"image":   "",
			"name":    "b",
			"role":    "",
			"score":   0,
		},
		{
			"country": "",
			"email":   "",
			"image":   "",
			"name":    "c",
			"role":    "",
			"score":   998,
		},
		{
			"country": "",
			"email":   "",
			"image":   "",
			"name":    "d",
			"role":    "",
			"score":   0,
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
			"country": "",
			"email":   "",
			"image":   "",
			"name":    "a",
			"role":    "",
			"score":   1498,
		},
		{
			"country": "",
			"email":   "",
			"image":   "",
			"name":    "b",
			"role":    "",
			"score":   0,
		},
		{
			"country": "",
			"email":   "",
			"image":   "",
			"name":    "c",
			"role":    "",
			"score":   998,
		},
		{
			"country": "",
			"email":   "",
			"image":   "",
			"name":    "d",
			"role":    "",
			"score":   0,
		},
		{
			"country": "",
			"email":   "",
			"image":   "",
			"name":    "test",
			"role":    "",
			"score":   0,
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
			"country": "",
			"email":   "a@a",
			"image":   "",
			"name":    "a",
			"role":    "Player",
			"score":   1498,
		},
		{
			"country": "",
			"email":   "b@b",
			"image":   "",
			"name":    "b",
			"role":    "Player",
			"score":   0,
		},
		{
			"country": "",
			"email":   "c@c",
			"image":   "",
			"name":    "c",
			"role":    "Player",
			"score":   998,
		},
		{
			"country": "",
			"email":   "d@d",
			"image":   "",
			"name":    "d",
			"role":    "Player",
			"score":   0,
		},
		{
			"country": "",
			"email":   "e@e",
			"image":   "",
			"name":    "e",
			"role":    "Admin",
			"score":   0,
		},
		{
			"country": "",
			"email":   "f@f",
			"image":   "",
			"name":    "f",
			"role":    "Author",
			"score":   0,
		},
		{
			"country": "",
			"email":   "test@test.test",
			"image":   "",
			"name":    "test",
			"role":    "Player",
			"score":   0,
		},
		{
			"country": "",
			"email":   "admin@test.com",
			"image":   "",
			"name":    "admin",
			"role":    "Admin",
			"score":   0,
		},
	}

	admin, err := user_register.RegisterUser(t.Context(), "admin", "admin@test.com", "testpass", sqlc.UserRoleAdmin)
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
