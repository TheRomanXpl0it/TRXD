package team_get_test

import (
	"fmt"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/user_register"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "team_get")
}

func TestTeamGet(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	A, err := db.GetTeamByName(t.Context(), "A")
	if err != nil {
		t.Fatalf("Failed to get team A: %v", err)
	}
	if A == nil {
		t.Fatal("Team A not found")
	}

	expectedPlayer := map[string]interface{}{
		"badges": []map[string]interface{}{
			{
				"description": "Completed all cat-1 challenges",
				"name":        "cat-1",
			},
		},
		"country": "",
		"id":      A.ID,
		"members": []map[string]interface{}{
			{
				"name":  "a",
				"role":  "Player",
				"score": 1498,
			},
			{
				"name":  "b",
				"role":  "Player",
				"score": 0,
			},
		},
		"name":  "A",
		"score": 1498,
		"solves": []map[string]interface{}{
			{
				"category": "cat-1",
				"name":     "chall-1",
			},
			{
				"category": "cat-1",
				"name":     "chall-3",
			},
			{
				"category": "cat-1",
				"name":     "chall-4",
			},
		},
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Get(fmt.Sprintf("/teams/%d", A.ID), nil, http.StatusOK)
	body := session.Body()
	for _, member := range body.(map[string]interface{})["members"].([]interface{}) {
		delete(member.(map[string]interface{}), "id")
	}
	for _, solve := range body.(map[string]interface{})["solves"].([]interface{}) {
		delete(solve.(map[string]interface{}), "id")
		delete(solve.(map[string]interface{}), "timestamp")
	}
	err = utils.Compare(expectedPlayer, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"username": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/%d", A.ID), nil, http.StatusOK)
	body = session.Body()
	for _, member := range body.(map[string]interface{})["members"].([]interface{}) {
		delete(member.(map[string]interface{}), "id")
	}
	for _, solve := range body.(map[string]interface{})["solves"].([]interface{}) {
		delete(solve.(map[string]interface{}), "id")
		delete(solve.(map[string]interface{}), "timestamp")
	}
	err = utils.Compare(expectedPlayer, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	expectedAdmin := map[string]interface{}{
		"badges": []map[string]interface{}{
			{
				"description": "Completed all cat-1 challenges",
				"name":        "cat-1",
			},
		},
		"country": "",
		"id":      A.ID,
		"members": []map[string]interface{}{
			{
				"name":  "a",
				"role":  "Player",
				"score": 1498,
			},
			{
				"name":  "b",
				"role":  "Player",
				"score": 0,
			},
			{
				"name":  "e",
				"role":  "Admin",
				"score": 0,
			},
		},
		"name":  "A",
		"score": 1498,
		"solves": []map[string]interface{}{
			{
				"category": "cat-1",
				"name":     "chall-1",
			},
			{
				"category": "cat-1",
				"name":     "chall-3",
			},
			{
				"category": "cat-1",
				"name":     "chall-4",
			},
		},
	}

	user, err := user_register.RegisterUser(t.Context(), "admin", "admin@admin.com", "adminpass", sqlc.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@admin.com", "password": "adminpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/teams/%d", A.ID), nil, http.StatusOK)
	body = session.Body()
	for _, member := range body.(map[string]interface{})["members"].([]interface{}) {
		delete(member.(map[string]interface{}), "id")
	}
	for _, solve := range body.(map[string]interface{})["solves"].([]interface{}) {
		delete(solve.(map[string]interface{}), "id")
		delete(solve.(map[string]interface{}), "timestamp")
	}
	err = utils.Compare(expectedAdmin, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}
}
