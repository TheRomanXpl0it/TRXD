package teams_scoreboard_graph_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/utils/test_utils"
)

type JSON map[string]any

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

	A := test_utils.GetTeamByName(t, "A")
	B := test_utils.GetTeamByName(t, "B")
	C := test_utils.GetTeamByName(t, "C")

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	var challID1, challID2, challID3, challID4 int32
	for _, chall := range List(body) {
		switch Json(chall)["name"] {
		case "chall-1":
			challID1 = Int32(Json(chall)["id"])
		case "chall-2":
			challID2 = Int32(Json(chall)["id"])
		case "chall-3":
			challID3 = Int32(Json(chall)["id"])
		case "chall-4":
			challID4 = Int32(Json(chall)["id"])
		}
	}

	expected := []JSON{
		{
			"submissions": []JSON{
				{
					"chall_id":    challID1,
					"first_blood": true,
					"score":       500,
				},
				{
					"chall_id":    challID3,
					"first_blood": true,
					"score":       1000,
				},
				{
					"chall_id":    challID4,
					"first_blood": true,
					"score":       1498,
				},
			},
			"team_id":   A.ID,
			"team_name": A.Name,
		},
		{
			"submissions": []JSON{
				{
					"chall_id":    challID4,
					"first_blood": false,
					"score":       498,
				},
				{
					"chall_id":    challID2,
					"first_blood": true,
					"score":       998,
				},
			},
			"team_id":   B.ID,
			"team_name": B.Name,
		},
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Get("/scoreboard/graph", nil, http.StatusOK)
	session.CheckFilteredResponse(expected, "timestamp")

	test_utils.UpdateConfig(t, "scoreboard-top", "3")
	session.Get("/scoreboard/graph", nil, http.StatusOK)
	session.CheckFilteredResponse(expected, "timestamp")

	test_utils.UpdateConfig(t, "scoreboard-top", "2")
	session.Get("/scoreboard/graph", nil, http.StatusOK)
	session.CheckFilteredResponse(expected, "timestamp")

	test_utils.UpdateConfig(t, "scoreboard-top", "1")
	session.Get("/scoreboard/graph", nil, http.StatusOK)
	session.CheckFilteredResponse(expected[:1], "timestamp")

	test_utils.UpdateConfig(t, "scoreboard-top", "0")
	session.Get("/scoreboard/graph", nil, http.StatusOK)
	session.CheckFilteredResponse([]JSON{}, "timestamp")

	test_utils.UpdateConfig(t, "scoreboard-top", "10")
	sessionC := test_utils.NewApiTestSession(t, app)
	sessionC.Post("/register", JSON{"name": "test-user", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	sessionC.Post("/teams/join", JSON{"name": "C", "password": "testpass"}, http.StatusOK)
	sessionC.Post("/submissions", JSON{"chall_id": challID1, "flag": "flag{test-1}"}, http.StatusOK)
	sessionC.Post("/submissions", JSON{"chall_id": challID2, "flag": "flag{test-2}"}, http.StatusOK)
	sessionC.Post("/submissions", JSON{"chall_id": challID3, "flag": "flag{test-3}"}, http.StatusOK)

	expected3 := []JSON{
		{
			"submissions": []JSON{
				{
					"chall_id":    challID1,
					"first_blood": true,
					"score":       498,
				},
				{
					"chall_id":    challID3,
					"first_blood": true,
					"score":       996,
				},
				{
					"chall_id":    challID4,
					"first_blood": true,
					"score":       1494,
				},
			},
			"team_id":   A.ID,
			"team_name": A.Name,
		},
		{
			"submissions": []JSON{
				{
					"chall_id":    challID1,
					"first_blood": false,
					"score":       498,
				},
				{
					"chall_id":    challID2,
					"first_blood": false,
					"score":       996,
				},
				{
					"chall_id":    challID3,
					"first_blood": false,
					"score":       1494,
				},
			},
			"team_id":   C.ID,
			"team_name": C.Name,
		},
		{
			"submissions": []JSON{
				{
					"chall_id":    challID4,
					"first_blood": false,
					"score":       498,
				},
				{
					"chall_id":    challID2,
					"first_blood": true,
					"score":       996,
				},
			},
			"team_id":   B.ID,
			"team_name": B.Name,
		},
	}

	session.Get("/scoreboard/graph", nil, http.StatusOK)
	session.CheckFilteredResponse(expected3, "timestamp")
}
