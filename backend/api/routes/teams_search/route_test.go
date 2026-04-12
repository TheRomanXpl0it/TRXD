package teams_search_test

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

	// SEARCH NAME

	playerName := "A"
	adminName := "admin"

	test_utils.RegisterUser(t, adminName, "admin@test.com", "testpass", sqlc.UserRoleAdmin)

	selfName := "self"
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": selfName, "email": "self@test.com", "password": "testpass"}, http.StatusOK)

	expectedNoAuth := JSON{
		"badges": []JSON{
			{
				"description": "Completed all cat-1 challenges",
				"name":        "cat-1",
			},
		},
		"country": "",
		"members": []JSON{
			{
				"name":  "a",
				"role":  "Player",
				"score": 1498,
			},
			{
				"name":  "b",
				"role":  "Player",
				"score": 0,
			},
		},
		"name":  playerName,
		"score": 1498,
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

	session.Get("/teams/search", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.MissingRequiredFields))

	session.Get("/teams/search?name=", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.MissingRequiredFields))

	session.Get("/teams/search?name=AAA", nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.TeamNotFound))

	session.Get("/teams/search?name="+strings.Repeat("A", consts.MaxUserNameLen+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(test_utils.Format(consts.MaxError, "team_name", consts.MaxUserNameLen)))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/teams/search?name=%s", adminName), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.TeamNotFound))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/teams/search?name=%s", playerName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedNoAuth, "id", "user_id", "timestamp")

	expectedPlayer := JSON{
		"badges": []JSON{
			{
				"description": "Completed all cat-1 challenges",
				"name":        "cat-1",
			},
		},
		"country": "",
		"members": []JSON{
			{
				"name":  "a",
				"role":  "Player",
				"score": 1498,
			},
			{
				"name":  "b",
				"role":  "Player",
				"score": 0,
			},
		},
		"name":  playerName,
		"score": 1498,
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
		"badges":  []JSON{},
		"country": "",
		"members": []JSON{
			{
				"name":  selfName,
				"role":  "Player",
				"score": 0,
			},
		},
		"name":   selfName,
		"score":  0,
		"solves": []JSON{},
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
	session.Get(fmt.Sprintf("/teams/search?name=%s", adminName), nil, http.StatusNotFound)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/search?name=%s", playerName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayer, "id", "user_id", "timestamp")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/search?name=%s", selfName), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.TeamNotFound))

	session.Post("/teams/register", JSON{"name": selfName, "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Get(fmt.Sprintf("/teams/search?name=%s", selfName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedSelf, "id", "user_id", "timestamp")

	expectedPlayerAdmin := JSON{
		"badges": []JSON{
			{
				"description": "Completed all cat-1 challenges",
				"name":        "cat-1",
			},
		},
		"country": "",
		"members": []JSON{
			{
				"name":  "a",
				"role":  "Player",
				"score": 1498,
			},
			{
				"name":  "b",
				"role":  "Player",
				"score": 0,
			},
			{
				"name":  "e",
				"role":  "Admin",
				"score": 0,
			},
		},
		"name":  "A",
		"score": 1498,
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
		"badges":  []JSON{},
		"country": "",
		"members": []JSON{
			{
				"name":  adminName,
				"role":  "Admin",
				"score": 0,
			},
		},
		"name":   adminName,
		"score":  0,
		"solves": []JSON{},
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
	session.Get(fmt.Sprintf("/teams/search?name=%s", playerName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayerAdmin, "id", "user_id", "timestamp")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/search?name=%s", adminName), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.TeamNotFound))

	session.Post("/teams/register", JSON{"name": adminName, "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Get(fmt.Sprintf("/teams/search?name=%s", adminName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedAdmin, "id", "user_id", "timestamp")

	test_utils.UpdateConfig(t, "start-time", time.Now().Add(10*time.Hour).Format(time.RFC3339))
	delete(expectedPlayerAdmin, "total_category_challenges")
	delete(expectedAdmin, "total_category_challenges")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/search?name=%s", playerName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayerAdmin, "id", "user_id", "timestamp")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.com", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/search?name=%s", adminName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedAdmin, "id", "user_id", "timestamp")

	test_utils.UpdateConfig(t, "start-time", "")

	// SEARCH EMAIL

	unregisteredEmail := "invalid@email.com"
	playerEmail := "a@a.a"
	adminEmail := "admin@test2.com"

	session = test_utils.NewApiTestSession(t, app)
	session.Get("/teams/search?email=AAA", nil, http.StatusUnauthorized)
	session.CheckResponse(errorf(consts.Unauthorized))

	session.Post("/login", JSON{"email": "self@test.com", "password": "testpass"}, http.StatusOK)
	session.Get("/teams/search?email=AAA", nil, http.StatusUnauthorized)
	session.CheckResponse(errorf(consts.Unauthorized))

	test_utils.RegisterUser(t, "admin2", adminEmail, "testpass", sqlc.UserRoleAdmin)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": adminEmail, "password": "testpass"}, http.StatusOK)

	session.Get("/teams/search", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.MissingRequiredFields))

	session.Get("/teams/search?email=", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.MissingRequiredFields))

	session.Get("/teams/search?email=AAA", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidEmail))

	session.Get("/teams/search?email="+strings.Repeat("A", consts.MaxEmailLen+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(test_utils.Format(consts.MaxError, "user_email", consts.MaxEmailLen)))

	session.Get(fmt.Sprintf("/teams/search?email=%s", unregisteredEmail), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.TeamNotFound))

	expectedPlayerAdmin = JSON{
		"badges": []JSON{
			{
				"description": "Completed all cat-1 challenges",
				"name":        "cat-1",
			},
		},
		"country": "",
		"members": []JSON{
			{
				"name":  "a",
				"role":  "Player",
				"score": 1498,
			},
			{
				"name":  "b",
				"role":  "Player",
				"score": 0,
			},
			{
				"name":  "e",
				"role":  "Admin",
				"score": 0,
			},
		},
		"name":  "A",
		"score": 1498,
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
	expectedAdmin = JSON{
		"badges":  []JSON{},
		"country": "",
		"members": []JSON{
			{
				"name":  "admin2",
				"role":  "Admin",
				"score": 0,
			},
		},
		"name":   "admin2",
		"score":  0,
		"solves": []JSON{},
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

	session.Get(fmt.Sprintf("/teams/search?email=%s", playerEmail), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayerAdmin, "id", "user_id", "timestamp")

	session.Get(fmt.Sprintf("/teams/search?email=%s", adminEmail), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.TeamNotFound))

	session.Post("/teams/register", JSON{"name": "admin2", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Get(fmt.Sprintf("/teams/search?email=%s", adminEmail), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedAdmin, "id", "user_id", "timestamp")

	test_utils.UpdateConfig(t, "start-time", time.Now().Add(10*time.Hour).Format(time.RFC3339))
	delete(expectedPlayerAdmin, "total_category_challenges")
	delete(expectedAdmin, "total_category_challenges")

	session.Get(fmt.Sprintf("/teams/search?email=%s", playerEmail), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayerAdmin, "id", "user_id", "timestamp")

	session.Get(fmt.Sprintf("/teams/search?email=%s", adminEmail), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedAdmin, "id", "user_id", "timestamp")

	test_utils.UpdateConfig(t, "start-time", "")

	// USER MODE

	test_utils.UpdateConfig(t, "user-mode", "true")

	userName := "user1"
	userEmail := "user1@example.com"

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": userName, "email": userEmail, "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)

	expected := JSON{
		"badges":  []JSON{},
		"country": "",
		"email":   userEmail,
		"name":    userName,
		"role":    "Player",
		"score":   0,
		"solves":  []JSON{},
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

	expectedPlayer = JSON{
		"badges": []JSON{
			{
				"description": "Completed all cat-1 challenges",
				"name":        "cat-1",
			},
		},
		"country": "",
		"name":    "A",
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

	session.Get(fmt.Sprintf("/teams/search?name=%s", userName), nil, http.StatusOK)
	session.CheckFilteredResponse(expected, "id", "user_id", "timestamp")

	session.Get(fmt.Sprintf("/teams/search?email=%s", userEmail), nil, http.StatusUnauthorized)
	session.CheckResponse(errorf(consts.Unauthorized))

	session.Get(fmt.Sprintf("/teams/search?name=%s", playerName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayer, "id", "user_id", "timestamp")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": adminEmail, "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/search?name=%s", userName), nil, http.StatusOK)
	session.CheckFilteredResponse(expected, "id", "user_id", "timestamp")

	session.Get(fmt.Sprintf("/teams/search?email=%s", userEmail), nil, http.StatusOK)
	session.CheckFilteredResponse(expected, "id", "user_id", "timestamp")

	expectedPlayer["email"] = playerEmail
	session.Get(fmt.Sprintf("/teams/search?name=%s", playerName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayer, "id", "user_id", "timestamp")

	session.Get(fmt.Sprintf("/teams/search?email=%s", playerEmail), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayer, "id", "user_id", "timestamp")

	expectedAdmin = JSON{
		"badges":  []JSON{},
		"country": "",
		"email":   "admin@test.com",
		"name":    "admin",
		"role":    "Admin",
		"score":   0,
		"solves":  []JSON{},
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

	session.Get(fmt.Sprintf("/teams/search?name=%s", adminName), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedAdmin, "id", "user_id", "timestamp")

	session.Get(fmt.Sprintf("/teams/search?email=%s", expectedAdmin["email"]), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedAdmin, "id", "user_id", "timestamp")
}
