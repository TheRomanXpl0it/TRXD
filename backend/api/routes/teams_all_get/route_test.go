package teams_all_get_test

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

	A := test_utils.GetTeamByName(t, "A")
	B := test_utils.GetTeamByName(t, "B")
	C := test_utils.GetTeamByName(t, "C")

	expected := JSON{
		"teams": []JSON{
			{
				"badges": []JSON{
					{
						"description": "Completed all cat-1 challenges",
						"name":        "cat-1",
					},
				},
				"country": "",
				"id":      A.ID,
				"name":    "A",
				"score":   1498,
			},
			{
				"badges": []JSON{
					{
						"description": "Completed all cat-2 challenges",
						"name":        "cat-2",
					},
				},
				"country": "",
				"id":      B.ID,
				"name":    "B",
				"score":   998,
			},
			{
				"badges":  []JSON{},
				"country": "",
				"id":      C.ID,
				"name":    "C",
				"score":   0,
			},
		},
		"total": 3,
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Get("/teams", nil, http.StatusOK)
	session.CheckResponse(expected)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/teams", nil, http.StatusOK)
	session.CheckResponse(expected)

	test_utils.RegisterUser(t, "admin", "admin@admin.com", "adminpass", sqlc.UserRoleAdmin)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@admin.com", "password": "adminpass"}, http.StatusOK)
	session.Get("/teams", nil, http.StatusOK)
	session.CheckResponse(expected)

	session.Get("/teams?offset=-1", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get("/teams?limit=-1", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get(fmt.Sprintf("/teams?offset=%d", math.MaxInt32+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get(fmt.Sprintf("/teams?limit=%d", math.MaxInt32+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))

	subSet := func(expected JSON, start int, end int) JSON {
		return JSON{
			"teams": expected["teams"].([]JSON)[start:end],
			"total": expected["total"],
		}
	}

	session.Get("/teams?offset=1&limit=1", nil, http.StatusOK)
	session.CheckResponse(subSet(expected, 1, 2))
	session.Get("/teams?offset=1&limit=2", nil, http.StatusOK)
	session.CheckResponse(subSet(expected, 1, 3))
	session.Get("/teams?offset=1", nil, http.StatusOK)
	session.CheckResponse(subSet(expected, 1, len(expected["teams"].([]JSON))))
	session.Get("/teams?limit=2", nil, http.StatusOK)
	session.CheckResponse(subSet(expected, 0, 2))
}
