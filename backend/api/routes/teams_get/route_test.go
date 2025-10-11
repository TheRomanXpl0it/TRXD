package teams_get_test

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

	session := test_utils.NewApiTestSession(t, app)
	session.Get("/teams/AAAA", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidTeamID))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/teams/%d", -1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf("team_id must be at least 0"))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/teams/%d", 99999), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.TeamNotFound))

	session = test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/teams/%d", A.ID), nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "timestamp", "user_id")
	test_utils.Compare(t, expectedPlayer, body)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/%d", A.ID), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "timestamp", "user_id")
	test_utils.Compare(t, expectedPlayer, body)

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

	test_utils.RegisterUser(t, "admin", "admin@admin.com", "adminpass", sqlc.UserRoleAdmin)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@admin.com", "password": "adminpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/%d", A.ID), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "timestamp", "user_id")
	test_utils.Compare(t, expectedAdmin, body)
}
