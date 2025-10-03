package users_get_test

import (
	"fmt"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	test_utils.RegisterUser(t, "admin", "admin@test.com", "testpass", sqlc.UserRoleAdmin)

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get("/users", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	idPlayer := int32(body.([]interface{})[0].(map[string]interface{})["id"].(float64))
	idAdmin := int32(body.([]interface{})[len(body.([]interface{}))-1].(map[string]interface{})["id"].(float64))

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "self", "email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get("/info", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	idSelf := int32(body.(map[string]interface{})["id"].(float64))

	expectedNoAuth := JSON{
		"country": "",
		"email":   "",
		"image":   "",
		"name":    "a",
		"role":    "",
		"score":   1498,
		"solves": []JSON{
			{
				"category": "cat-1",
				"name":     "chall-1",
			},
			{
				"category": "cat-1",
				"name":     "chall-3",
			},
			{
				"category": "cat-1",
				"name":     "chall-4",
			},
		},
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Get("/users/AAA", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidUserID))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/%d", -1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf("user_id must be at least 0"))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/%d", 99999), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.UserNotFound))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/%d", idAdmin), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.UserNotFound))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/%d", idPlayer), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "joined_at", "team_id", "timestamp")
	test_utils.Compare(t, expectedNoAuth, body)

	expectedPlayer := JSON{
		"country": "",
		"email":   "",
		"image":   "",
		"name":    "a",
		"role":    "",
		"score":   1498,
		"solves": []JSON{
			{
				"category": "cat-1",
				"name":     "chall-1",
			},
			{
				"category": "cat-1",
				"name":     "chall-3",
			},
			{
				"category": "cat-1",
				"name":     "chall-4",
			},
		},
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
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idAdmin), nil, http.StatusNotFound)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idPlayer), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "joined_at", "team_id", "timestamp")
	test_utils.Compare(t, expectedPlayer, body)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idSelf), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "joined_at", "team_id", "timestamp")
	test_utils.Compare(t, expectedSelf, body)

	expectedPlayerAdmin := JSON{
		"country": "",
		"email":   "a@a.a",
		"image":   "",
		"name":    "a",
		"role":    "Player",
		"score":   1498,
		"solves": []JSON{
			{
				"category": "cat-1",
				"name":     "chall-1",
			},
			{
				"category": "cat-1",
				"name":     "chall-3",
			},
			{
				"category": "cat-1",
				"name":     "chall-4",
			},
		},
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
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idPlayer), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "joined_at", "team_id", "timestamp")
	test_utils.Compare(t, expectedPlayerAdmin, body)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idAdmin), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "joined_at", "timestamp")
	test_utils.Compare(t, expectedAdmin, body)
}
