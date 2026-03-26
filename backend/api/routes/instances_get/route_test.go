package instances_get_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils/test_utils"
)

type JSON map[string]any

func Json(val any) map[string]any {
	return val.(map[string]any)
}

func List(val any) []any {
	return val.([]any)
}

func Int32(val any) int32 {
	return int32(val.(float64))
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	test_utils.RegisterUser(t, "test", "test@test.test", "testpass", sqlc.UserRolePlayer)

	admin := test_utils.NewApiTestSession(t, app)
	admin.Post("/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)
	admin.CheckResponse(nil)

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams/register", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()

	var challID3, challID4 int32
	for _, chall := range List(body) {
		switch Json(chall)["name"] {
		case "chall-3":
			challID3 = Int32(Json(chall)["id"])
		case "chall-4":
			challID4 = Int32(Json(chall)["id"])
		}
	}

	expected := []JSON{}
	admin.Get("/instances", nil, http.StatusOK)
	admin.CheckResponse(expected)

	expected = []JSON{
		{
			"chall_id":   challID3,
			"chall_name": "chall-3",
			"conn_type":  "HTTP",
			"port":       0,
			"team_name":  "test-team",
		},
		{
			"chall_id":   challID4,
			"chall_name": "chall-4",
			"conn_type":  "HTTP",
			"port":       0,
			"team_name":  "test-team",
		},
	}
	session.Post("/instances", JSON{"chall_id": challID3}, http.StatusOK)
	session.Post("/instances", JSON{"chall_id": challID4}, http.StatusOK)
	admin.Get("/instances", nil, http.StatusOK)
	admin.CheckFilteredResponse(expected, "expires_at", "host", "team_id", "docker_id")

	expected = []JSON{
		{
			"chall_id":   challID3,
			"chall_name": "chall-3",
			"conn_type":  "HTTP",
			"port":       0,
			"team_name":  "test-team",
		},
		{
			"chall_id":   challID4,
			"chall_name": "chall-4",
			"conn_type":  "HTTP",
			"port":       0,
			"team_name":  "test-team",
		},
		{
			"chall_id":   challID3,
			"chall_name": "chall-3",
			"conn_type":  "HTTP",
			"port":       0,
			"team_name":  "A",
		},
		{
			"chall_id":   challID4,
			"chall_name": "chall-4",
			"conn_type":  "HTTP",
			"port":       0,
			"team_name":  "A",
		},
	}
	admin.Post("/instances", JSON{"chall_id": challID3}, http.StatusOK)
	admin.Post("/instances", JSON{"chall_id": challID4}, http.StatusOK)
	admin.Get("/instances", nil, http.StatusOK)
	admin.CheckFilteredResponse(expected, "expires_at", "host", "team_id", "docker_id")

	expected = []JSON{
		{
			"chall_id":   challID3,
			"chall_name": "chall-3",
			"conn_type":  "HTTP",
			"port":       0,
			"team_name":  "test-team",
		},
		{
			"chall_id":   challID4,
			"chall_name": "chall-4",
			"conn_type":  "HTTP",
			"port":       0,
			"team_name":  "A",
		},
	}
	session.Delete("/instances", JSON{"chall_id": challID4}, http.StatusOK)
	admin.Delete("/instances", JSON{"chall_id": challID3}, http.StatusOK)
	admin.Get("/instances", nil, http.StatusOK)
	admin.CheckFilteredResponse(expected, "expires_at", "host", "team_id", "docker_id")

	expected = []JSON{}
	session.Delete("/instances", JSON{"chall_id": challID3}, http.StatusOK)
	admin.Delete("/instances", JSON{"chall_id": challID4}, http.StatusOK)
	admin.Get("/instances", nil, http.StatusOK)
	admin.CheckResponse(expected)
}
