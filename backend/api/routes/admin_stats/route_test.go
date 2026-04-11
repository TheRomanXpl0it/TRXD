package admin_stats_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/utils/test_utils"
)

type JSON map[string]any

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	expected := JSON{
		"total_challenges":          5,
		"total_correct_submissions": 6,
		"total_players":             4,
		"total_released_challenges": 4,
		"total_submissions":         12,
		"total_teams":               3,
		"total_users":               6,
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Get("/stats", nil, http.StatusOK)
	session.CheckResponse(expected)
}
