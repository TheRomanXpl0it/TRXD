package challenges_all_get_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "challenges_all_get")
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	expectedPlayer := []JSON{
		{
			"category":   "cat-1",
			"difficulty": "Easy",
			"hidden":     false,
			"instance":   false,
			"name":       "chall-1",
			"points":     500,
			"solved":     false,
			"solves":     1,
			"tags": []interface{}{
				"tag-1",
				"test-tag",
			},
		},
		{
			"category":   "cat-1",
			"difficulty": "Hard",
			"hidden":     false,
			"instance":   true,
			"name":       "chall-3",
			"points":     500,
			"solved":     false,
			"solves":     1,
			"tags": []interface{}{
				"tag-3",
			},
		},
		{
			"category":   "cat-1",
			"difficulty": "Insane",
			"hidden":     false,
			"instance":   true,
			"name":       "chall-4",
			"points":     498,
			"solved":     false,
			"solves":     2,
			"tags": []interface{}{
				"tag-4",
			},
		},
		{
			"category":   "cat-2",
			"difficulty": "Medium",
			"hidden":     false,
			"instance":   false,
			"name":       "chall-2",
			"points":     500,
			"solved":     false,
			"solves":     1,
			"tags": []interface{}{
				"tag-2",
			},
		},
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/users/register", JSON{"username": "test", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams/register", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	for _, chall := range body.([]interface{}) {
		delete(chall.(map[string]interface{}), "id")
	}
	err := utils.Compare(expectedPlayer, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	expectedAuthor := []JSON{
		{
			"category":   "cat-2",
			"difficulty": "Easy",
			"hidden":     true,
			"instance":   false,
			"name":       "chall-5",
			"points":     500,
			"solved":     false,
			"solves":     0,
			"tags": []interface{}{
				"tag-5",
			},
		},
		{
			"category":   "cat-1",
			"difficulty": "Easy",
			"hidden":     false,
			"instance":   false,
			"name":       "chall-1",
			"points":     500,
			"solved":     false,
			"solves":     1,
			"tags": []interface{}{
				"tag-1",
				"test-tag",
			},
		},
		{
			"category":   "cat-1",
			"difficulty": "Hard",
			"hidden":     false,
			"instance":   true,
			"name":       "chall-3",
			"points":     500,
			"solved":     false,
			"solves":     1,
			"tags": []interface{}{
				"tag-3",
			},
		},
		{
			"category":   "cat-1",
			"difficulty": "Insane",
			"hidden":     false,
			"instance":   true,
			"name":       "chall-4",
			"points":     498,
			"solved":     false,
			"solves":     2,
			"tags": []interface{}{
				"tag-4",
			},
		},
		{
			"category":   "cat-2",
			"difficulty": "Medium",
			"hidden":     false,
			"instance":   false,
			"name":       "chall-2",
			"points":     500,
			"solved":     false,
			"solves":     1,
			"tags": []interface{}{
				"tag-2",
			},
		},
	}

	test_utils.RegisterUser(t, "test2", "test3@test.test", "testpass", sqlc.UserRoleAuthor)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "test3@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body = session.Body()
	for _, chall := range body.([]interface{}) {
		delete(chall.(map[string]interface{}), "id")
	}
	err = utils.Compare(expectedAuthor, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}
}
