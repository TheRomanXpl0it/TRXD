package challenges_get_test

import (
	"context"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/user_register"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "challenges_submit")
}

func TestGetChallenges(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	expectedPlayer := []JSON{
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
				"author2",
			},
			"category":    "cat-1",
			"description": "TEST chall-1 DESC",
			"difficulty":  "Easy",
			"flags":       nil,
			"hidden":      false,
			"host":        "http://theromanxpl0.it",
			"instance":    false,
			"name":        "chall-1",
			"points":      500,
			"port":        1337,
			"solved":      false,
			"solves":      1,
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
			},
			"tags": []interface{}{
				"tag-1",
				"test-tag",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
			},
			"category":    "cat-1",
			"description": "TEST chall-3 DESC",
			"difficulty":  "Hard",
			"flags":       nil,
			"hidden":      false,
			"host":        "",
			"instance":    true,
			"name":        "chall-3",
			"points":      500,
			"port":        0,
			"solved":      false,
			"solves":      1,
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
			},
			"tags": []interface{}{
				"tag-3",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author2",
			},
			"category":    "cat-1",
			"description": "TEST chall-4 DESC",
			"difficulty":  "Insane",
			"flags":       nil,
			"hidden":      false,
			"host":        "",
			"instance":    true,
			"name":        "chall-4",
			"points":      498,
			"port":        0,
			"solved":      false,
			"solves":      2,
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
				{
					"name": "B",
				},
			},
			"tags": []interface{}{
				"tag-4",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
				"author2",
				"author3",
			},
			"category":    "cat-2",
			"description": "TEST chall-2 DESC",
			"difficulty":  "Medium",
			"flags":       nil,
			"hidden":      false,
			"host":        "",
			"instance":    false,
			"name":        "chall-2",
			"points":      500,
			"port":        0,
			"solved":      false,
			"solves":      1,
			"solves_list": []map[string]interface{}{
				{
					"name": "B",
				},
			},
			"tags": []interface{}{
				"tag-2",
			},
			"timeout": 0,
		},
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"username": "test", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	for _, chall := range body.([]interface{}) {
		delete(chall.(map[string]interface{}), "id")
		for _, solve := range chall.(map[string]interface{})["solves_list"].([]interface{}) {
			delete(solve.(map[string]interface{}), "id")
			delete(solve.(map[string]interface{}), "timestamp")
		}
	}
	err := utils.Compare(expectedPlayer, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	expectedAuthor := []JSON{
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author3",
			},
			"category":    "cat-2",
			"description": "TEST chall-5 DESC",
			"difficulty":  "Easy",
			"flags": []map[string]interface{}{
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
			"solves_list": []map[string]interface{}{},
			"tags": []interface{}{
				"tag-5",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
				"author2",
			},
			"category":    "cat-1",
			"description": "TEST chall-1 DESC",
			"difficulty":  "Easy",
			"flags": []map[string]interface{}{
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
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
			},
			"tags": []interface{}{
				"tag-1",
				"test-tag",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
			},
			"category":    "cat-1",
			"description": "TEST chall-3 DESC",
			"difficulty":  "Hard",
			"flags": []map[string]interface{}{
				{
					"flag":  "flag{test-3}",
					"regex": false,
				},
			},
			"hidden":   false,
			"host":     "",
			"instance": true,
			"name":     "chall-3",
			"points":   500,
			"port":     0,
			"solved":   false,
			"solves":   1,
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
			},
			"tags": []interface{}{
				"tag-3",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author2",
			},
			"category":    "cat-1",
			"description": "TEST chall-4 DESC",
			"difficulty":  "Insane",
			"flags": []map[string]interface{}{
				{
					"flag":  "flag{test-4}",
					"regex": false,
				},
			},
			"hidden":   false,
			"host":     "",
			"instance": true,
			"name":     "chall-4",
			"points":   498,
			"port":     0,
			"solved":   false,
			"solves":   2,
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
				{
					"name": "B",
				},
			},
			"tags": []interface{}{
				"tag-4",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
				"author2",
				"author3",
			},
			"category":    "cat-2",
			"description": "TEST chall-2 DESC",
			"difficulty":  "Medium",
			"flags": []map[string]interface{}{
				{
					"flag":  "flag{test-2}",
					"regex": false,
				},
			},
			"hidden":   false,
			"host":     "",
			"instance": false,
			"name":     "chall-2",
			"points":   500,
			"port":     0,
			"solved":   false,
			"solves":   1,
			"solves_list": []map[string]interface{}{
				{
					"name": "B",
				},
			},
			"tags": []interface{}{
				"tag-2",
			},
			"timeout": 0,
		},
	}

	user, err := user_register.RegisterUser(context.Background(), "test2", "test3@test.test", "testpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test3@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body = session.Body()
	for _, chall := range body.([]interface{}) {
		delete(chall.(map[string]interface{}), "id")
		for _, solve := range chall.(map[string]interface{})["solves_list"].([]interface{}) {
			delete(solve.(map[string]interface{}), "id")
			delete(solve.(map[string]interface{}), "timestamp")
		}
	}
	err = utils.Compare(expectedAuthor, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}
}
