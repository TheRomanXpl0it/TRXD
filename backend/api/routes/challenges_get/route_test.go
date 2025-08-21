package challenges_get_test

import (
	"fmt"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/users/register", JSON{"username": "test", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams/register", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	id := int32(body.([]interface{})[0].(map[string]interface{})["id"].(float64))

	expectedPlayer := JSON{
		"attachments": []interface{}{},
		"authors": []interface{}{
			"author1",
			"author2",
		},
		"category":    "cat-1",
		"description": "TEST chall-1 DESC",
		"difficulty":  "Easy",
		"first_blood": JSON{
			"name": "A",
		},
		"host":     "http://theromanxpl0.it",
		"instance": false,
		"name":     "chall-1",
		"points":   500,
		"port":     1337,
		"solved":   false,
		"solves":   1,
		"solves_list": []JSON{
			{
				"name": "A",
			},
		},
		"tags": []interface{}{
			"tag-1",
			"test-tag",
		},
		"timeout": 0,
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/challenges/%d", id), nil, http.StatusOK)
	body = session.Body()
	test_utils.DeleteKeys(body, "id", "timestamp")
	test_utils.Compare(t, expectedPlayer, body)

	expectedAuthor := JSON{
		"attachments": []interface{}{},
		"authors": []interface{}{
			"author1",
			"author2",
		},
		"category":    "cat-1",
		"description": "TEST chall-1 DESC",
		"difficulty":  "Easy",
		"first_blood": JSON{
			"name": "A",
		},
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
		"hidden":     false,
		"host":       "http://theromanxpl0.it",
		"instance":   false,
		"max_points": 500,
		"name":       "chall-1",
		"points":     500,
		"port":       1337,
		"score_type": "Dynamic",
		"solved":     false,
		"solves":     1,
		"solves_list": []JSON{
			{
				"name": "A",
			},
		},
		"tags": []interface{}{
			"tag-1",
			"test-tag",
		},
		"timeout": 0,
		"type":    "Normal",
	}

	test_utils.RegisterUser(t, "test2", "test3@test.test", "testpass", sqlc.UserRoleAuthor)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "test3@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/challenges/%d", id), nil, http.StatusOK)
	body = session.Body()
	test_utils.DeleteKeys(body, "id", "timestamp")
	test_utils.Compare(t, expectedAuthor, body)

	expectedAuthorHidden := JSON{
		"attachments": []interface{}{},
		"authors": []interface{}{
			"author3",
		},
		"category":    "cat-2",
		"description": "TEST chall-5 DESC",
		"difficulty":  "Easy",
		"first_blood": nil,
		"flags": []JSON{
			{
				"flag":  "flag{test-5}",
				"regex": false,
			},
		},
		"hidden":      true,
		"host":        "",
		"instance":    false,
		"max_points":  500,
		"name":        "chall-5",
		"points":      500,
		"port":        0,
		"score_type":  "Static",
		"solved":      false,
		"solves":      0,
		"solves_list": []JSON{},
		"tags": []interface{}{
			"tag-5",
		},
		"timeout": 0,
		"type":    "Normal",
	}

	session.Get("/challenges", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	id = int32(body.([]interface{})[0].(map[string]interface{})["id"].(float64))

	session.Get(fmt.Sprintf("/challenges/%d", id), nil, http.StatusOK)
	body = session.Body()
	test_utils.DeleteKeys(body, "id", "timestamp")
	test_utils.Compare(t, expectedAuthorHidden, body)

	// TODO: add test with dockeconfig
	// TODO: add test with instances
}
