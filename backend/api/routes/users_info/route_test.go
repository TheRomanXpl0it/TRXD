package users_info_test

import (
	"net/http"
	"testing"
	"time"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	session := test_utils.NewApiTestSession(t, app)

	session.Get("/info", nil, http.StatusOK)
	session.CheckResponse(nil)

	session.Post("/register", JSON{"name": "test", "email": "allow@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)

	expected := JSON{
		"name":      "test",
		"role":      sqlc.UserRolePlayer,
		"team_id":   nil,
		"user_mode": false,
	}
	session.Get("/info", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, expected, body)

	startTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	endTime := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
	test_utils.UpdateConfig(t, "start-time", startTime)
	test_utils.UpdateConfig(t, "end-time", endTime)

	expected = JSON{
		"end_time":   endTime,
		"name":       "test",
		"role":       sqlc.UserRolePlayer,
		"start_time": startTime,
		"team_id":    nil,
		"user_mode":  false,
	}
	session.Get("/info", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, expected, body)

	test_utils.UpdateConfig(t, "start-time", "")
	test_utils.UpdateConfig(t, "end-time", "")
	session.Post("/teams/register", JSON{"name": "test", "password": "testpass"}, http.StatusOK)

	expected = JSON{
		"name":      "test",
		"role":      sqlc.UserRolePlayer,
		"user_mode": false,
	}
	session.Get("/info", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	test_utils.DeleteKeys(body, "id")
	if body.(map[string]interface{})["team_id"] == nil {
		t.Errorf("Expected team_id to be set, got nil")
	}
	test_utils.DeleteKeys(body, "team_id")
	test_utils.Compare(t, expected, body)
}
