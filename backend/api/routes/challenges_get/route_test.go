package challenges_get_test

import (
	"fmt"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "challenges_get")
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
			"id":   1,
			"name": "A",
		},
		"flags":    nil,
		"hidden":   false,
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
	delete(body.(map[string]interface{}), "id")
	for _, solve := range body.(map[string]interface{})["solves_list"].([]interface{}) {
		delete(solve.(map[string]interface{}), "id")
		delete(solve.(map[string]interface{}), "timestamp")
	}
	err := utils.Compare(expectedPlayer, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

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
			"id":   1,
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
		"hidden":   false,
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

	test_utils.RegisterUser(t, "test2", "test3@test.test", "testpass", sqlc.UserRoleAuthor)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "test3@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/challenges/%d", id), nil, http.StatusOK)
	body = session.Body()
	delete(body.(map[string]interface{}), "id")
	for _, solve := range body.(map[string]interface{})["solves_list"].([]interface{}) {
		delete(solve.(map[string]interface{}), "id")
		delete(solve.(map[string]interface{}), "timestamp")
	}
	err = utils.Compare(expectedAuthor, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

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
		"name":        "chall-5",
		"points":      500,
		"port":        0,
		"solved":      false,
		"solves":      0,
		"solves_list": []JSON{},
		"tags": []interface{}{
			"tag-5",
		},
		"timeout": 0,
	}

	session.Get("/challenges", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	id = int32(body.([]interface{})[0].(map[string]interface{})["id"].(float64))

	session.Get(fmt.Sprintf("/challenges/%d", id), nil, http.StatusOK)
	body = session.Body()
	delete(body.(map[string]interface{}), "id")
	for _, solve := range body.(map[string]interface{})["solves_list"].([]interface{}) {
		delete(solve.(map[string]interface{}), "id")
		delete(solve.(map[string]interface{}), "timestamp")
	}
	err = utils.Compare(expectedAuthorHidden, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}
}
