package users_get_test

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

type JSON map[string]any

func errorf(val any) JSON {
	return JSON{"error": val}
}

func Json(val any) map[string]any {
	return val.(map[string]any)
}

func List(val any) []any {
	return val.([]any)
}

func Int32(val any) int32 {
	return int32(val.(float64))
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	test_utils.RegisterUser(t, "admin", "admin@test.com", "testpass", sqlc.UserRoleAdmin)

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get("/users", nil, http.StatusOK)
	body := session.Body()
	var idPlayer, idAdmin int32
	for _, user := range List(Json(body)["users"]) {
		switch Json(user)["name"] {
		case "a":
			idPlayer = Int32(Json(user)["id"])
		case "admin":
			idAdmin = Int32(Json(user)["id"])
		}
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "self", "email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get("/info", nil, http.StatusOK)
	body = session.Body()
	idSelf := Int32(Json(body)["id"])

	expectedNoAuth := JSON{
		"country": "",
		"name":    "a",
		"score":   1498,
		"solves": []JSON{
			{
				"category":    "cat-1",
				"first_blood": true,
				"name":        "chall-1",
				"points":      500,
			},
			{
				"category":    "cat-1",
				"first_blood": true,
				"name":        "chall-3",
				"points":      500,
			},
			{
				"category":    "cat-1",
				"first_blood": true,
				"name":        "chall-4",
				"points":      498,
			},
		},
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Get("/users/AAA", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidUserID))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/%d", -1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(test_utils.Format(consts.MinError, "id", 0)))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/%d", 99999), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.UserNotFound))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/%d", math.MaxInt32+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(test_utils.Format(consts.MinError, "id", 0)))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/%d", idAdmin), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.UserNotFound))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/%d", idPlayer), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedNoAuth, "id", "joined_at", "team_id", "timestamp")

	expectedPlayer := JSON{
		"country": "",
		"name":    "a",
		"score":   1498,
		"solves": []JSON{
			{
				"category":    "cat-1",
				"first_blood": true,
				"name":        "chall-1",
				"points":      500,
			},
			{
				"category":    "cat-1",
				"first_blood": true,
				"name":        "chall-3",
				"points":      500,
			},
			{
				"category":    "cat-1",
				"first_blood": true,
				"name":        "chall-4",
				"points":      498,
			},
		},
	}
	expectedSelf := JSON{
		"country": "",
		"email":   "self@test.com",
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
	session.CheckFilteredResponse(expectedPlayer, "id", "joined_at", "team_id", "timestamp")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idSelf), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedSelf, "id", "joined_at", "team_id", "timestamp")

	expectedPlayerAdmin := JSON{
		"country": "",
		"email":   "a@a.a",
		"name":    "a",
		"role":    "Player",
		"score":   1498,
		"solves": []JSON{
			{
				"category":    "cat-1",
				"first_blood": true,
				"name":        "chall-1",
				"points":      500,
			},
			{
				"category":    "cat-1",
				"first_blood": true,
				"name":        "chall-3",
				"points":      500,
			},
			{
				"category":    "cat-1",
				"first_blood": true,
				"name":        "chall-4",
				"points":      498,
			},
		},
	}
	expectedAdmin := JSON{
		"country": "",
		"email":   "admin@test.com",
		"name":    "admin",
		"role":    "Admin",
		"score":   0,
		"team_id": nil,
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idPlayer), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayerAdmin, "id", "joined_at", "team_id", "timestamp")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/%d", idAdmin), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedAdmin, "id", "joined_at", "timestamp")
}
