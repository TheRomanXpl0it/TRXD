package users_get_name_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
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

	playerName := "a"
	adminName := "admin"

	test_utils.RegisterUser(t, adminName, "admin@test.com", "testpass", sqlc.UserRoleAdmin)

	selfName := "self"
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": selfName, "email": "self@test.com", "password": "testpass"}, http.StatusOK)

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
		"total_category_challenges": []JSON{
			{
				"category": "cat-1",
				"count":    3,
			},
			{
				"category": "cat-2",
				"count":    1,
			},
		},
	}

	session = test_utils.NewApiTestSession(t, app)

	session.Get("/users/name", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidUserName))

	session.Get("/users/name?name=", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidUserName))

	session.Get("/users/name?name=AAA", nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.UserNotFound))

	session.Get("/users/name?name="+strings.Repeat("A", consts.MaxUserNameLen+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(test_utils.Format(consts.MaxError, "user_name", consts.MaxUserNameLen)))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/name?name=%s", adminName), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.UserNotFound))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/users/name?name=%s", playerName), nil, http.StatusOK)
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
		"total_category_challenges": []JSON{
			{
				"category": "cat-1",
				"count":    3,
			},
			{
				"category": "cat-2",
				"count":    1,
			},
		},
	}
	expectedSelf := JSON{
		"country": "",
		"email":   "self@test.com",
		"name":    "self",
		"role":    "Player",
		"score":   0,
		"total_category_challenges": []JSON{
			{
				"category": "cat-1",
				"count":    3,
			},
			{
				"category": "cat-2",
				"count":    1,
			},
		},
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/name?name=%s", adminName), nil, http.StatusNotFound)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/name?name=%s", playerName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayer, "id", "joined_at", "team_id", "timestamp")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/name?name=%s", selfName), nil, http.StatusOK)
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
		"total_category_challenges": []JSON{
			{
				"category": "cat-1",
				"count":    3,
			},
			{
				"category": "cat-2",
				"count":    1,
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
		"total_category_challenges": []JSON{
			{
				"category": "cat-1",
				"count":    3,
			},
			{
				"category": "cat-2",
				"count":    1,
			},
		},
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/name?name=%s", playerName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayerAdmin, "id", "joined_at", "team_id", "timestamp")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/name?name=%s", adminName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedAdmin, "id", "joined_at", "timestamp")

	test_utils.UpdateConfig(t, "start-time", time.Now().Add(10*time.Hour).Format(time.RFC3339))
	delete(expectedPlayerAdmin, "total_category_challenges")
	delete(expectedAdmin, "total_category_challenges")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/name?name=%s", playerName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayerAdmin, "id", "joined_at", "team_id", "timestamp")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/users/name?name=%s", adminName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedAdmin, "id", "joined_at", "timestamp")

}
