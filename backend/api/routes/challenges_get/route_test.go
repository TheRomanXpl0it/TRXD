package challenges_get_test

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
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams/register", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	id := int32(body.([]interface{})[0].(map[string]interface{})["id"].(float64))

	expectedPlayer := JSON{
		"solves_list": []JSON{
			{
				"name": "A",
			},
		},
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges/AAA", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidChallengeID))

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/challenges/%d", -1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(test_utils.Format(consts.MinError, "id", 0)))

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/challenges/%d", 99999), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.ChallengeNotFound))

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/challenges/%d", math.MaxInt32+1), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(test_utils.Format(consts.MinError, "id", 0)))

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/challenges/%d", id), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "timestamp")
	test_utils.Compare(t, expectedPlayer, body)

	expectedAuthor := JSON{
		"flags": []JSON{
			{
				"flag":  "flag{test-1}",
				"regex": false,
			},
			{
				"flag":  "flag\\{test-[a-z]{2}\\}",
				"regex": true,
			},
		},
		"solves_list": []JSON{
			{
				"name": "A",
			},
		},
		"type": "Normal",
	}

	test_utils.RegisterUser(t, "test2", "test3@test.test", "testpass", sqlc.UserRoleAuthor)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test3@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams/register", JSON{"name": "test-team-2", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/challenges/%d", id), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "timestamp")
	test_utils.Compare(t, expectedAuthor, body)

	expectedAuthorHidden := JSON{
		"flags": []JSON{
			{
				"flag":  "flag{test-5}",
				"regex": false,
			},
		},
		"solves_list": []interface{}{},
		"type":        "Normal",
	}

	session.Get("/challenges", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	id5 := int32(body.([]interface{})[4].(map[string]interface{})["id"].(float64))
	id3 := int32(body.([]interface{})[2].(map[string]interface{})["id"].(float64))

	session.Get(fmt.Sprintf("/challenges/%d", id5), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "timestamp")
	test_utils.Compare(t, expectedAuthorHidden, body)

	expectedDocker := JSON{
		"docker_config": JSON{
			"compose":     "",
			"envs":        "",
			"hash_domain": true,
			"image":       "echo-server:latest",
			"lifetime":    0,
			"max_cpu":     "",
			"max_memory":  0,
		},
		"flags": []JSON{
			{
				"flag":  "flag{test-3}",
				"regex": false,
			},
		},
		"solves_list": []JSON{
			{
				"name": "A",
			},
		},
		"type": "Container",
	}

	session.Get(fmt.Sprintf("/challenges/%d", id3), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "timestamp")
	test_utils.Compare(t, expectedDocker, body)

	session.Post("/instances", JSON{"chall_id": id3}, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}

	expectedInstance := JSON{
		"docker_config": JSON{
			"compose":     "",
			"envs":        "",
			"hash_domain": true,
			"image":       "echo-server:latest",
			"lifetime":    0,
			"max_cpu":     "",
			"max_memory":  0,
		},
		"flags": []JSON{
			{
				"flag":  "flag{test-3}",
				"regex": false,
			},
		},
		"solves_list": []JSON{
			{
				"name": "A",
			},
		},
		"type": "Container",
	}

	session.Get(fmt.Sprintf("/challenges/%d", id3), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id", "timestamp", "timeout")
	test_utils.Compare(t, expectedInstance, body)
}
