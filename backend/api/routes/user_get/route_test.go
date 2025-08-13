package user_get_test

import (
	"context"
	"fmt"
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

func TestUserGet(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	admin, err := user_register.RegisterUser(context.Background(), "admin", "admin@test.com", "testpass", db.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}
	if admin == nil {
		t.Fatal("Admin registration returned nil")
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get("/users", nil, http.StatusOK)
	body := session.Body()
	idPlayer := int32(body.([]interface{})[0].(map[string]interface{})["id"].(float64))
	idAdmin := int32(body.([]interface{})[len(body.([]interface{}))-1].(map[string]interface{})["id"].(float64))

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"username": "self", "email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get("/info", nil, http.StatusOK)
	body = session.Body()
	idSelf := int32(body.(map[string]interface{})["id"].(float64))

	expectedNoAuth := JSON{
		"email":       "",
		"image":       "",
		"name":        "a",
		"nationality": "",
		"role":        "",
		"score":       1498,
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/%d", idAdmin), nil, http.StatusNotFound)

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/%d", idPlayer), nil, http.StatusOK)
	body = session.Body()
	delete(body.(map[string]interface{}), "id")
	delete(body.(map[string]interface{}), "joined_at")
	delete(body.(map[string]interface{}), "solves")
	delete(body.(map[string]interface{}), "team_id")
	err = utils.Compare(expectedNoAuth, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	expectedPlayer := JSON{
		"email":       "",
		"image":       "",
		"name":        "a",
		"nationality": "",
		"role":        "",
		"score":       1498,
	}
	expectedSelf := JSON{
		"email":       "self@test.com",
		"image":       "",
		"name":        "self",
		"nationality": "",
		"role":        "Player",
		"score":       0,
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idAdmin), nil, http.StatusNotFound)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idPlayer), nil, http.StatusOK)
	body = session.Body()
	delete(body.(map[string]interface{}), "id")
	delete(body.(map[string]interface{}), "joined_at")
	delete(body.(map[string]interface{}), "solves")
	delete(body.(map[string]interface{}), "team_id")
	err = utils.Compare(expectedPlayer, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idSelf), nil, http.StatusOK)
	body = session.Body()
	delete(body.(map[string]interface{}), "id")
	delete(body.(map[string]interface{}), "joined_at")
	delete(body.(map[string]interface{}), "solves")
	delete(body.(map[string]interface{}), "team_id")
	err = utils.Compare(expectedSelf, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	expectedPlayerAdmin := JSON{
		"email":       "a@a",
		"image":       "",
		"name":        "a",
		"nationality": "",
		"role":        "Player",
		"score":       1498,
	}
	expectedAdmin := JSON{
		"email":       "admin@test.com",
		"image":       "",
		"name":        "admin",
		"nationality": "",
		"role":        "Admin",
		"score":       0,
		"team_id":     nil,
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idPlayer), nil, http.StatusOK)
	body = session.Body()
	delete(body.(map[string]interface{}), "id")
	delete(body.(map[string]interface{}), "joined_at")
	delete(body.(map[string]interface{}), "solves")
	delete(body.(map[string]interface{}), "team_id")
	err = utils.Compare(expectedPlayerAdmin, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idAdmin), nil, http.StatusOK)
	body = session.Body()
	delete(body.(map[string]interface{}), "id")
	delete(body.(map[string]interface{}), "joined_at")
	delete(body.(map[string]interface{}), "solves")
	err = utils.Compare(expectedAdmin, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}
}
