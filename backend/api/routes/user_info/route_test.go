package user_info_test

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
	test_utils.Main(m, "../../../", "user_info")
}

func TestUserInfo(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)

	session.Get("/info", nil, http.StatusUnauthorized)
	session.CheckResponse(errorf(consts.Unauthorized))

	session.Post("/register", JSON{"username": "test", "email": "allow@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)

	session.Get("/info", nil, http.StatusOK)
	body := session.Body().(map[string]interface{})
	delete(body, "id")
	utils.Compare(body, JSON{"username": "test", "role": sqlc.UserRolePlayer, "team_id": nil})

	session.Post("/teams", JSON{"name": "test", "password": "testpass"}, http.StatusOK)

	session.Get("/info", nil, http.StatusOK)
	body = session.Body().(map[string]interface{})
	delete(body, "id")
	if body["team_id"] == nil {
		t.Errorf("Expected team_id to be set, got nil")
	}
	delete(body, "team_id")
	utils.Compare(body, JSON{"username": "test", "role": sqlc.UserRolePlayer})
}
