package instances_update_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
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

	test_utils.RegisterUser(t, "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)
	test_utils.RegisterUser(t, "test", "test@test.test", "testpass", sqlc.UserRolePlayer)

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Post("/teams/register", JSON{"name": "author-team", "password": "authorpass"}, http.StatusOK)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}

	var challID1, challID3, challID4 int32
	for _, chall := range body.([]interface{}) {
		switch chall.(map[string]interface{})["name"] {
		case "chall-1":
			challID1 = int32(chall.(map[string]interface{})["id"].(float64))
		case "chall-3":
			challID3 = int32(chall.(map[string]interface{})["id"].(float64))
		case "chall-4":
			challID4 = int32(chall.(map[string]interface{})["id"].(float64))
		}
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Patch("/instances", nil, http.StatusForbidden)

	session.Post("/teams/register", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
	session.Patch("/instances", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidJSON))

	session.Patch("/instances", JSON{}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.MissingRequiredFields))

	session.Patch("/instances", JSON{"chall_id": -1}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidChallengeID))

	session.Patch("/instances", JSON{"chall_id": 99999}, http.StatusNotFound)
	session.CheckResponse(errorf(consts.ChallengeNotFound))

	test_utils.UpdateConfig(t, "secret", "")
	session.Patch("/instances", JSON{"chall_id": challID1}, http.StatusForbidden)
	session.CheckResponse(errorf(consts.DisabledInstances))
	test_utils.UpdateConfig(t, "secret", "test-secret")

	session.Patch("/instances", JSON{"chall_id": challID1}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.ChallengeNotInstanciable))

	session.Delete("/instances", JSON{"chall_id": challID4}, http.StatusNotFound)
	session.CheckResponse(errorf(consts.InstanceNotFound))

	session.Patch("/instances", JSON{"chall_id": challID3}, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	if _, ok := body.(map[string]interface{})["timeout"]; !ok {
		t.Fatalf("Expected timeout to be present in response: %+v", body)
	}
}
