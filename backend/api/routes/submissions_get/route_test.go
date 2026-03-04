package submissions_get_test

import (
	"fmt"
	"math"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]any

func errorf(val any) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	expected := JSON{
		"submissions": []JSON{
			{
				"chall_name":  "chall-3",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "C",
				"user_name":   "f",
			},
			{
				"chall_name":  "chall-2",
				"first_blood": true,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "B",
				"user_name":   "c",
			},
			{
				"chall_name":  "chall-4",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "B",
				"user_name":   "c",
			},
			{
				"chall_name":  "chall-3",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Repeated",
				"team_name":   "A",
				"user_name":   "b",
			},
			{
				"chall_name":  "chall-1",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Repeated",
				"team_name":   "A",
				"user_name":   "a",
			},
			{
				"chall_name":  "chall-4",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Repeated",
				"team_name":   "A",
				"user_name":   "a",
			},
			{
				"chall_name":  "chall-4",
				"first_blood": true,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "A",
				"user_name":   "a",
			},
			{
				"chall_name":  "chall-3",
				"first_blood": true,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "A",
				"user_name":   "a",
			},
			{
				"chall_name":  "chall-1",
				"first_blood": true,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "A",
				"user_name":   "a",
			},
			{
				"chall_name":  "chall-1",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Repeated",
				"team_name":   "A",
				"user_name":   "a",
			},
			{
				"chall_name":  "chall-1",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Wrong",
				"team_name":   "A",
				"user_name":   "a",
			},
		},
		"total": 12,
	}
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)
	session.Get("/submissions", nil, http.StatusOK)
	session.CheckFilteredResponse(expected, "id", "user_id", "team_id", "chall_id", "timestamp")

	session.Get("/submissions?offset=-1", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get("/submissions?limit=-1", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get(fmt.Sprintf("/submissions?offset=%d", math.MaxInt32+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get(fmt.Sprintf("/submissions?limit=%d", math.MaxInt32+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))

	subSet := func(expected JSON, start int, end int) JSON {
		return JSON{
			"submissions": expected["submissions"].([]JSON)[start:end],
			"total":       expected["total"],
		}
	}

	session.Get("/submissions?offset=1", nil, http.StatusOK)
	sub := subSet(expected, 1, len(expected["submissions"].([]JSON)))
	session.CheckFilteredResponse(sub, "id", "user_id", "team_id", "chall_id", "timestamp")

	session.Get("/submissions?limit=2", nil, http.StatusOK)
	sub = subSet(expected, 0, 2)
	session.CheckFilteredResponse(sub, "id", "user_id", "team_id", "chall_id", "timestamp")

	session.Get("/submissions?offset=1&limit=2", nil, http.StatusOK)
	sub = subSet(expected, 1, 3)
	session.CheckFilteredResponse(sub, "id", "user_id", "team_id", "chall_id", "timestamp")

	// User Mode

	test_utils.UpdateConfig(t, "user-mode", "true")
	expected = JSON{
		"submissions": []JSON{
			{
				"chall_name":  "chall-3",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "C",
			},
			{
				"chall_name":  "chall-2",
				"first_blood": true,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "B",
			},
			{
				"chall_name":  "chall-4",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "B",
			},
			{
				"chall_name":  "chall-3",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Repeated",
				"team_name":   "A",
			},
			{
				"chall_name":  "chall-1",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Repeated",
				"team_name":   "A",
			},
			{
				"chall_name":  "chall-4",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Repeated",
				"team_name":   "A",
			},
			{
				"chall_name":  "chall-4",
				"first_blood": true,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "A",
			},
			{
				"chall_name":  "chall-3",
				"first_blood": true,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "A",
			},
			{
				"chall_name":  "chall-1",
				"first_blood": true,
				"flag":        "flag",
				"status":      "Correct",
				"team_name":   "A",
			},
			{
				"chall_name":  "chall-1",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Repeated",
				"team_name":   "A",
			},
			{
				"chall_name":  "chall-1",
				"first_blood": false,
				"flag":        "flag",
				"status":      "Wrong",
				"team_name":   "A",
			},
		},
		"total": 12,
	}
	session.Get("/submissions", nil, http.StatusOK)
	session.CheckFilteredResponse(expected, "id", "team_id", "chall_id", "timestamp")
}
