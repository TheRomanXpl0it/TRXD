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
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
				"author2",
			},
			"category":    "cat-1",
			"description": "TEST chall-1 DESC",
			"difficulty":  "Easy",
			"first_blood": false,
			"hidden":      false,
			"host":        "http://theromanxpl0.it",
			"instance":    false,
			"max_points":  500,
			"name":        "chall-1",
			"points":      500,
			"port":        1234,
			"score_type":  "Dynamic",
			"solved":      false,
			"solves":      1,
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
			"first_blood": false,
			"hidden":      false,
			"host":        "chall-3.test.com",
			"instance":    true,
			"max_points":  500,
			"name":        "chall-3",
			"points":      500,
			"port":        1337,
			"score_type":  "Dynamic",
			"solved":      false,
			"solves":      1,
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
			"first_blood": false,
			"hidden":      false,
			"host":        "",
			"instance":    true,
			"max_points":  500,
			"name":        "chall-4",
			"points":      498,
			"port":        0,
			"score_type":  "Dynamic",
			"solved":      false,
			"solves":      2,
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
			"first_blood": false,
			"hidden":      false,
			"host":        "",
			"instance":    false,
			"max_points":  500,
			"name":        "chall-2",
			"points":      500,
			"port":        0,
			"score_type":  "Dynamic",
			"solved":      false,
			"solves":      1,
			"tags": []interface{}{
				"tag-2",
			},
			"timeout": 0,
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
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author3",
			},
			"category":    "cat-2",
			"description": "TEST chall-5 DESC",
			"difficulty":  "Easy",
			"first_blood": false,
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
			"tags": []interface{}{
				"tag-5",
			},
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
			"first_blood": false,
			"hidden":      false,
			"host":        "http://theromanxpl0.it",
			"instance":    false,
			"max_points":  500,
			"name":        "chall-1",
			"points":      500,
			"port":        1234,
			"score_type":  "Dynamic",
			"solved":      true,
			"solves":      1,
			"tags": []interface{}{
				"tag-1",
				"test-tag",
			},
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
			},
			"category":    "cat-1",
			"description": "TEST chall-3 DESC",
			"difficulty":  "Hard",
			"first_blood": false,
			"hidden":      false,
			"host":        "f6735eddbadf.chall-3.test.com",
			"instance":    true,
			"max_points":  500,
			"name":        "chall-3",
			"points":      500,
			"port":        1337,
			"score_type":  "Dynamic",
			"solved":      true,
			"solves":      1,
			"tags": []interface{}{
				"tag-3",
			},
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author2",
			},
			"category":    "cat-1",
			"description": "TEST chall-4 DESC",
			"difficulty":  "Insane",
			"first_blood": false,
			"hidden":      false,
			"host":        "",
			"instance":    true,
			"max_points":  500,
			"name":        "chall-4",
			"points":      498,
			"port":        0,
			"score_type":  "Dynamic",
			"solved":      true,
			"solves":      2,
			"tags": []interface{}{
				"tag-4",
			},
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
			"first_blood": false,
			"hidden":      false,
			"host":        "",
			"instance":    false,
			"max_points":  500,
			"name":        "chall-2",
			"points":      500,
			"port":        0,
			"score_type":  "Dynamic",
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

	// TODO: instance test
}
