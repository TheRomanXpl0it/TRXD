package teams_get_test

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

func Int32(val any) int32 {
	return int32(val.(float64))
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	A := test_utils.GetTeamByName(t, "A")

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
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Get("/teams/AAAA", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidTeamID))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/teams/%d", -1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(test_utils.Format(consts.MinError, "id", 0)))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/teams/%d", 99999), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.TeamNotFound))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/teams/%d", math.MaxInt32+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(test_utils.Format(consts.MinError, "id", 0)))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/teams/%d", A.ID), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayer, "id", "timestamp", "user_id")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/%d", A.ID), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedPlayer, "id", "timestamp", "user_id")

	expectedAdmin := JSON{
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
	}

	test_utils.RegisterUser(t, "admin", "admin@admin.com", "adminpass", sqlc.UserRoleAdmin)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@admin.com", "password": "adminpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/%d", A.ID), nil, http.StatusOK)
	session.CheckFilteredResponse(expectedAdmin, "id", "timestamp", "user_id")

	//! User Mode tests
	test_utils.UpdateConfig(t, "user-mode", "true")
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "single", "email": "single@gmail.com", "password": "testpass"}, http.StatusOK)
	session.Get("/info", nil, http.StatusOK)
	body := session.Body()
	teamID := Int32(Json(body)["team_id"])
	userID := Int32(Json(body)["id"])

	expected := JSON{
		"badges":  []any{},
		"country": "",
		"email":   "single@gmail.com",
		"id":      teamID,
		"name":    "single",
		"role":    "Player",
		"score":   0,
		"solves":  []any{},
		"user_id": userID,
	}
	session.Get(fmt.Sprintf("/teams/%d", teamID), nil, http.StatusOK)
	session.CheckResponse(expected)

	expected = JSON{
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
	}
	session.Get(fmt.Sprintf("/teams/%d", A.ID), nil, http.StatusOK)
	session.CheckFilteredResponse(expected, "id", "timestamp", "user_id")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)

	expected["email"] = "a@a.a"
	session.Get(fmt.Sprintf("/teams/%d", A.ID), nil, http.StatusOK)
	session.CheckFilteredResponse(expected, "id", "timestamp", "user_id")

	expected = JSON{
		"badges":  []any{},
		"country": "",
		"email":   "single@gmail.com",
		"id":      teamID,
		"name":    "single",
		"role":    "Player",
		"score":   0,
		"solves":  []any{},
		"user_id": userID,
	}
	session.Get(fmt.Sprintf("/teams/%d", teamID), nil, http.StatusOK)
	session.CheckResponse(expected)
}
