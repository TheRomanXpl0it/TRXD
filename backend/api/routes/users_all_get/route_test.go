package users_all_get_test

import (
	"fmt"
	"math"
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
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	expectedNoAuth := JSON{
		"total": 4,
		"users": []JSON{
			{
				"country": "",
				"email":   "",
				"name":    "a",
				"role":    "",
				"score":   1498,
			},
			{
				"country": "",
				"email":   "",
				"name":    "b",
				"role":    "",
				"score":   0,
			},
			{
				"country": "",
				"email":   "",
				"name":    "c",
				"role":    "",
				"score":   998,
			},
			{
				"country": "",
				"email":   "",
				"name":    "d",
				"role":    "",
				"score":   0,
			},
		},
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Get("/users", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, expectedNoAuth, body)

	expectedPlayer := JSON{
		"total": 5,
		"users": []JSON{
			{
				"country": "",
				"email":   "",
				"name":    "a",
				"role":    "",
				"score":   1498,
			},
			{
				"country": "",
				"email":   "",
				"name":    "b",
				"role":    "",
				"score":   0,
			},
			{
				"country": "",
				"email":   "",
				"name":    "c",
				"role":    "",
				"score":   998,
			},
			{
				"country": "",
				"email":   "",
				"name":    "d",
				"role":    "",
				"score":   0,
			},
			{
				"country": "",
				"email":   "",
				"name":    "test",
				"role":    "",
				"score":   0,
			},
		},
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/users", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, expectedPlayer, body)

	session.Get("/users?offset=-1", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get("/users?limit=-1", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get(fmt.Sprintf("/users?offset=%d", math.MaxInt32+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get(fmt.Sprintf("/users?limit=%d", math.MaxInt32+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))

	subSet := func(expected JSON, start int, end int) JSON {
		return JSON{
			"users": expected["users"].([]JSON)[start:end],
			"total": expected["total"],
		}
	}

	session.Get("/users?offset=1", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, subSet(expectedPlayer, 1, len(expectedPlayer["users"].([]JSON))), body)

	session.Get("/users?limit=2", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, subSet(expectedPlayer, 0, 2), body)

	session.Get("/users?offset=1&limit=2", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, subSet(expectedPlayer, 1, 3), body)

	expectedAdmin := JSON{
		"total": 8,
		"users": []JSON{
			{
				"country": "",
				"email":   "a@a.a",
				"name":    "a",
				"role":    "Player",
				"score":   1498,
			},
			{
				"country": "",
				"email":   "b@b.b",
				"name":    "b",
				"role":    "Player",
				"score":   0,
			},
			{
				"country": "",
				"email":   "c@c.c",
				"name":    "c",
				"role":    "Player",
				"score":   998,
			},
			{
				"country": "",
				"email":   "d@d.d",
				"name":    "d",
				"role":    "Player",
				"score":   0,
			},
			{
				"country": "",
				"email":   "admin@email.com",
				"name":    "e",
				"role":    "Admin",
				"score":   0,
			},
			{
				"country": "",
				"email":   "f@f.f",
				"name":    "f",
				"role":    "Author",
				"score":   0,
			},
			{
				"country": "",
				"email":   "test@test.test",
				"name":    "test",
				"role":    "Player",
				"score":   0,
			},
			{
				"country": "",
				"email":   "admin@test.com",
				"name":    "admin",
				"role":    "Admin",
				"score":   0,
			},
		},
	}

	test_utils.RegisterUser(t, "admin", "admin@test.com", "testpass", sqlc.UserRoleAdmin)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get("/users", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, expectedAdmin, body)

	session.Get("/users?offset=1", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, subSet(expectedAdmin, 1, len(expectedAdmin["users"].([]JSON))), body)

	session.Get("/users?limit=2", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, subSet(expectedAdmin, 0, 2), body)

	session.Get("/users?offset=1&limit=2", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, subSet(expectedAdmin, 1, 3), body)

}
