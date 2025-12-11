package instances_create_test

import (
	"math"
	"net/http"
	"strings"
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
	app := api.SetupApp(t.Context())
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

	var challID1, challID2, challID3, challID4 int32
	for _, chall := range body.([]interface{}) {
		switch chall.(map[string]interface{})["name"] {
		case "chall-1":
			challID1 = int32(chall.(map[string]interface{})["id"].(float64))
		case "chall-2":
			challID2 = int32(chall.(map[string]interface{})["id"].(float64))
		case "chall-3":
			challID3 = int32(chall.(map[string]interface{})["id"].(float64))
		case "chall-4":
			challID4 = int32(chall.(map[string]interface{})["id"].(float64))
		}
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/instances", nil, http.StatusForbidden)

	session.Post("/teams/register", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
	session.Post("/instances", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidJSON))

	session.Post("/instances", JSON{}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.MissingRequiredFields))

	session.Post("/instances", JSON{"chall_id": -1}, http.StatusBadRequest)
	session.CheckResponse(errorf(test_utils.Format(consts.MinError, "ChallID", 0)))

	session.Post("/instances", JSON{"chall_id": 99999}, http.StatusNotFound)
	session.CheckResponse(errorf(consts.ChallengeNotFound))

	session.Post("/instances", JSON{"chall_id": math.MaxInt32 + 1}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidJSON))

	test_utils.UpdateConfig(t, "secret", "")
	session.Post("/instances", JSON{"chall_id": challID1}, http.StatusForbidden)
	session.CheckResponse(errorf(consts.DisabledInstances))
	test_utils.UpdateConfig(t, "secret", "test-secret")

	session.Post("/instances", JSON{"chall_id": challID1}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.ChallengeNotInstanciable))

	session.Post("/instances", JSON{"chall_id": challID3}, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	if host, ok := body.(map[string]interface{})["host"]; !ok {
		t.Fatalf("Expected host to be present in response: %+v", body)
	} else {
		if !strings.HasSuffix(host.(string), ".chall-3.test.com") {
			t.Fatalf("Expected host to end with .chall-3.test.com: %s", host)
		}
	}
	if _, ok := body.(map[string]interface{})["port"]; ok {
		t.Fatalf("Expected port to not be present in response: %+v", body)
	}
	if _, ok := body.(map[string]interface{})["timeout"]; !ok {
		t.Fatalf("Expected timeout to be present in response: %+v", body)
	}

	session.Post("/instances", JSON{"chall_id": challID3}, http.StatusConflict)
	session.CheckResponse(errorf(consts.AlreadyAnActiveInstance))

	session.Delete("/instances", JSON{"chall_id": challID3}, http.StatusOK)
	session.Delete("/instances", JSON{"chall_id": challID3}, http.StatusNotFound)
	session.CheckResponse(errorf(consts.InstanceNotFound))

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Patch("/challenges", JSON{"chall_id": challID3, "hash_domain": false}, http.StatusOK)
	session.CheckResponse(nil)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/instances", JSON{"chall_id": challID3}, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	if host, ok := body.(map[string]interface{})["host"]; !ok {
		t.Fatalf("Expected host to be present in response: %+v", body)
	} else {
		if host.(string) != "chall-3.test.com" {
			t.Fatalf("Expected host to be chall-3.test.com: %s", host)
		}
	}
	if _, ok := body.(map[string]interface{})["port"]; !ok {
		t.Fatalf("Expected port to be present in response: %+v", body)
	}
	if _, ok := body.(map[string]interface{})["timeout"]; !ok {
		t.Fatalf("Expected timeout to be present in response: %+v", body)
	}
	session.Delete("/instances", JSON{"chall_id": challID3}, http.StatusOK)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Patch("/challenges", JSON{"chall_id": challID3, "host": "", "hash_domain": true}, http.StatusOK)
	session.CheckResponse(nil)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/instances", JSON{"chall_id": challID3}, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	if host, ok := body.(map[string]interface{})["host"]; !ok {
		t.Fatalf("Expected host to be present in response: %+v", body)
	} else {
		if !strings.HasSuffix(host.(string), ".test.com") || strings.HasSuffix(host.(string), ".chall-3.test.com") {
			t.Fatalf("Expected host to end with .test.com: %s", host)
		}
	}
	if _, ok := body.(map[string]interface{})["port"]; ok {
		t.Fatalf("Expected port to not be present in response: %+v", body)
	}
	if _, ok := body.(map[string]interface{})["timeout"]; !ok {
		t.Fatalf("Expected timeout to be present in response: %+v", body)
	}
	session.Delete("/instances", JSON{"chall_id": challID3}, http.StatusOK)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Patch("/challenges", JSON{"chall_id": challID3, "type": "Compose"}, http.StatusOK)
	session.CheckResponse(nil)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/instances", JSON{"chall_id": challID3}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidImage))

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Patch("/challenges", JSON{"chall_id": challID3, "type": "Container"}, http.StatusOK)
	session.CheckResponse(nil)

	session.Patch("/challenges", JSON{"chall_id": challID2, "type": "Container", "image": "aaaa"}, http.StatusOK)
	session.CheckResponse(nil)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/instances", JSON{"chall_id": challID2}, http.StatusInternalServerError)
	session.CheckResponse(errorf(consts.ErrorCreatingInstance))

	session.Post("/instances", JSON{"chall_id": challID4}, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	if _, ok := body.(map[string]interface{})["host"]; !ok {
		t.Fatalf("Expected host to be present in response: %+v", body)
	}
	if _, ok := body.(map[string]interface{})["port"]; ok {
		t.Fatalf("Expected port to not be present in response: %+v", body)
	}
	if _, ok := body.(map[string]interface{})["timeout"]; !ok {
		t.Fatalf("Expected timeout to be present in response: %+v", body)
	}
}
