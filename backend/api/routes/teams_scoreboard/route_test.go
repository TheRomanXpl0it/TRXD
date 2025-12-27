package teams_scoreboard_test

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

	expected := []JSON{
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
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Get("/scoreboard", nil, http.StatusOK)
	session.CheckResponse(expected)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/scoreboard", nil, http.StatusOK)
	session.CheckResponse(expected)

	test_utils.RegisterUser(t, "player", "player@email.com", "testpass", sqlc.UserRolePlayer)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "player@email.com", "password": "testpass"}, http.StatusOK)
	session.Post("/teams/join", JSON{"name": "C", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	var challID1, challID2, challID3, challID4 int32
	for _, chall := range body.([]interface{}) {
		switch chall.(map[string]interface{})["name"] {
		case "chall-1":
			challID1 = int32(chall.(map[string]interface{})["id"].(float64))
		case "chall-2":
			challID2 = int32(chall.(map[string]interface{})["id"].(float64))
		case "chall-3":
			challID3 = int32(chall.(map[string]interface{})["id"].(float64))
		case "chall-4":
			challID4 = int32(chall.(map[string]interface{})["id"].(float64))
		}
	}
	session.Post("/submissions", JSON{"chall_id": challID1, "flag": "flag{test-1}"}, http.StatusOK)
	session.Post("/submissions", JSON{"chall_id": challID3, "flag": "flag{test-3}"}, http.StatusOK)
	session.Post("/submissions", JSON{"chall_id": challID4, "flag": "flag{test-4}"}, http.StatusOK)

	expected = []JSON{
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
			"score":   1488,
		},
		{
			"badges": []JSON{
				{
					"description": "Completed all cat-1 challenges",
					"name":        "cat-1",
				},
			},
			"country": "",
			"id":      C.ID,
			"name":    "C",
			"score":   1488,
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
			"score":   992,
		},
	}
	session.Get("/scoreboard", nil, http.StatusOK)
	session.CheckResponse(expected)

	session.Post("/submissions", JSON{"chall_id": challID2, "flag": "flag{test-2}"}, http.StatusOK)
	expected = []JSON{
		{
			"badges": []JSON{
				{
					"description": "Completed all cat-1 challenges",
					"name":        "cat-1",
				},
				{
					"description": "Completed all cat-2 challenges",
					"name":        "cat-2",
				},
			},
			"country": "",
			"id":      C.ID,
			"name":    "C",
			"score":   1986,
		},
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
			"score":   1488,
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
			"score":   990,
		},
	}
	session.Get("/scoreboard", nil, http.StatusOK)
	session.CheckResponse(expected)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "a@a.a", "password": "testpass"}, http.StatusOK)
	session.Post("/submissions", JSON{"chall_id": challID2, "flag": "flag{test-2}"}, http.StatusOK)
	expected = []JSON{
		{
			"badges": []JSON{
				{
					"description": "Completed all cat-1 challenges",
					"name":        "cat-1",
				},
				{
					"description": "Completed all cat-2 challenges",
					"name":        "cat-2",
				},
			},
			"country": "",
			"id":      C.ID,
			"name":    "C",
			"score":   1980,
		},
		{
			"badges": []JSON{
				{
					"description": "Completed all cat-1 challenges",
					"name":        "cat-1",
				},
				{
					"description": "Completed all cat-2 challenges",
					"name":        "cat-2",
				},
			},
			"country": "",
			"id":      A.ID,
			"name":    "A",
			"score":   1980,
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
			"score":   984,
		},
	}
	session.Get("/scoreboard", nil, http.StatusOK)
	session.CheckResponse(expected)

	test_utils.RegisterUser(t, "admin", "admin@test.test", "adminpass", sqlc.UserRoleAdmin)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.test", "password": "adminpass"}, http.StatusOK)
	session.Get("/scoreboard", nil, http.StatusOK)
	session.CheckResponse(expected)

	session.Get("/scoreboard?start=-1", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get("/scoreboard?end=-1", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get(fmt.Sprintf("/scoreboard?start=%d", math.MaxInt32+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get(fmt.Sprintf("/scoreboard?end=%d", math.MaxInt32+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))
	session.Get("/scoreboard?start=2&end=1", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidParam))

	session.Get("/scoreboard?start=1&end=2", nil, http.StatusOK)
	session.CheckResponse(expected[1:2])

	session.Get("/scoreboard?start=1", nil, http.StatusOK)
	session.CheckResponse(expected[1:])
}
