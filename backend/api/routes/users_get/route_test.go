package users_get_test

import (
	"fmt"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "users_get")
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	test_utils.RegisterUser(t, "admin", "admin@test.com", "testpass", sqlc.UserRoleAdmin)

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get("/users", nil, http.StatusOK)
	body := session.Body()
	idPlayer := int32(body.([]interface{})[0].(map[string]interface{})["id"].(float64))
	idAdmin := int32(body.([]interface{})[len(body.([]interface{}))-1].(map[string]interface{})["id"].(float64))

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/register", JSON{"username": "self", "email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get("/users/info", nil, http.StatusOK)
	body = session.Body()
	idSelf := int32(body.(map[string]interface{})["id"].(float64))

	expectedNoAuth := JSON{
		"country": "",
		"email":   "",
		"image":   "",
		"name":    "a",
		"role":    "",
		"score":   1498,
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
	err := utils.Compare(expectedNoAuth, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	expectedPlayer := JSON{
		"country": "",
		"email":   "",
		"image":   "",
		"name":    "a",
		"role":    "",
		"score":   1498,
	}
	expectedSelf := JSON{
		"country": "",
		"email":   "self@test.com",
		"image":   "",
		"name":    "self",
		"role":    "Player",
		"score":   0,
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idAdmin), nil, http.StatusNotFound)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
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
	session.Post("/users/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
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
		"country": "",
		"email":   "a@a",
		"image":   "",
		"name":    "a",
		"role":    "Player",
		"score":   1498,
	}
	expectedAdmin := JSON{
		"country": "",
		"email":   "admin@test.com",
		"image":   "",
		"name":    "admin",
		"role":    "Admin",
		"score":   0,
		"team_id": nil,
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
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
	session.Post("/users/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
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
