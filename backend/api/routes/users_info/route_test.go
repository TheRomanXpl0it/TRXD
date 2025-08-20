package users_info_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)

	session.Get("/users/info", nil, http.StatusUnauthorized)
	session.CheckResponse(errorf(consts.Unauthorized))

	session.Post("/users/register", JSON{"username": "test", "email": "allow@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)

	session.Get("/users/info", nil, http.StatusOK)
	body := session.Body()
	test_utils.DeleteKeys(body, "id")
	utils.Compare(body, JSON{"username": "test", "role": sqlc.UserRolePlayer, "team_id": nil})

	session.Post("/teams/register", JSON{"name": "test", "password": "testpass"}, http.StatusOK)

	session.Get("/users/info", nil, http.StatusOK)
	body = session.Body()
	test_utils.DeleteKeys(body, "id")
	if body.(map[string]interface{})["team_id"] == nil {
		t.Errorf("Expected team_id to be set, got nil")
	}
	test_utils.DeleteKeys(body, "team_id")
	test_utils.Compare(t, body, JSON{"username": "test", "role": sqlc.UserRolePlayer})
}
