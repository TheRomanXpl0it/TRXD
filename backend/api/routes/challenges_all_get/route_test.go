package challenges_all_get_test

import (
	"fmt"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	expectedPlayer := []JSON{
		{
			"category":    "cat-1",
			"difficulty":  "Easy",
			"first_blood": false,
			"hidden":      false,
			"instance":    false,
			"name":        "chall-1",
			"points":      500,
			"solved":      false,
			"solves":      1,
			"tags": []interface{}{
				"tag-1",
				"test-tag",
			},
		},
		{
			"category":    "cat-1",
			"difficulty":  "Hard",
			"first_blood": false,
			"hidden":      false,
			"instance":    true,
			"name":        "chall-3",
			"points":      500,
			"solved":      false,
			"solves":      1,
			"tags": []interface{}{
				"tag-3",
			},
		},
		{
			"category":    "cat-1",
			"difficulty":  "Insane",
			"first_blood": false,
			"hidden":      false,
			"instance":    true,
			"name":        "chall-4",
			"points":      498,
			"solved":      false,
			"solves":      2,
			"tags": []interface{}{
				"tag-4",
			},
		},
		{
			"category":    "cat-2",
			"difficulty":  "Medium",
			"first_blood": false,
			"hidden":      false,
			"instance":    false,
			"name":        "chall-2",
			"points":      500,
			"solved":      false,
			"solves":      1,
			"tags": []interface{}{
				"tag-2",
			},
		},
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams/register", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	var challID int32
	for _, chall := range body.([]interface{}) {
		if chall.(map[string]interface{})["name"] == "chall-3" {
			challID = int32(chall.(map[string]interface{})["id"].(float64))
			break
		}
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, expectedPlayer, body)

	expectedAuthor := []JSON{
		{
			"category":    "cat-2",
			"difficulty":  "Easy",
			"first_blood": false,
			"hidden":      true,
			"instance":    false,
			"name":        "chall-5",
			"points":      500,
			"solved":      false,
			"solves":      0,
			"tags": []interface{}{
				"tag-5",
			},
		},
		{
			"category":    "cat-1",
			"difficulty":  "Easy",
			"first_blood": true,
			"hidden":      false,
			"instance":    false,
			"name":        "chall-1",
			"points":      500,
			"solved":      true,
			"solves":      1,
			"tags": []interface{}{
				"tag-1",
				"test-tag",
			},
		},
		{
			"category":    "cat-1",
			"difficulty":  "Hard",
			"first_blood": true,
			"hidden":      false,
			"instance":    true,
			"name":        "chall-3",
			"points":      500,
			"solved":      true,
			"solves":      1,
			"tags": []interface{}{
				"tag-3",
			},
		},
		{
			"category":    "cat-1",
			"difficulty":  "Insane",
			"first_blood": true,
			"hidden":      false,
			"instance":    true,
			"name":        "chall-4",
			"points":      498,
			"solved":      true,
			"solves":      2,
			"tags": []interface{}{
				"tag-4",
			},
		},
		{
			"category":    "cat-2",
			"difficulty":  "Medium",
			"first_blood": false,
			"hidden":      false,
			"instance":    false,
			"name":        "chall-2",
			"points":      500,
			"solved":      false,
			"solves":      1,
			"tags": []interface{}{
				"tag-2",
			},
		},
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)
	session.Post("/instances", JSON{"chall_id": challID}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	conf, err := db.GetConfig(t.Context(), "instance-lifetime")
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	if conf == "" {
		t.Fatal("Expected config to not be nil")
	}
	var lifetime int
	_, err = fmt.Sscanf(conf, "%d", &lifetime)
	if err != nil {
		t.Fatalf("Failed to parse config value: %v", err)
	}
	if timeout, ok := body.([]interface{})[2].(map[string]interface{})["timeout"]; ok &&
		(int(timeout.(float64)) < lifetime-100 || int(timeout.(float64)) > lifetime) {
		t.Fatalf("Expected timeout to be around %d, got %v", lifetime, timeout)
	}
	test_utils.DeleteKeys(body, "id", "timeout")
	test_utils.Compare(t, expectedAuthor, body)
}
