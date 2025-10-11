package tags_delete_test

import (
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

var testData = []struct {
	testBody         interface{}
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"chall_id": ""},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": "", "name": strings.Repeat("a", consts.MaxTagNameLen+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf("Name must not exceed 32"),
	},
	{
		testBody:         JSON{"chall_id": -1, "name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf("ChallID must be at least 0"),
	},
	{
		testBody:       JSON{"chall_id": "", "name": "test-2"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"chall_id": "", "name": "test"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	test_utils.RegisterUser(t, "author", "author@test.test", "testpass", sqlc.UserRoleAuthor)
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/categories", JSON{"name": "cat", "icon": "icon"}, http.StatusOK)
	chall := test_utils.CreateChallenge(t, "chall", "cat", "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeStatic)

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "testpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = chall.ID
			}
		}
		session.Post("/tags", JSON{"chall_id": chall.ID, "name": "test"}, -1)
		session.Delete("/tags", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	var challengeBody interface{}
	for _, v := range body.([]interface{}) {
		if int32(v.(map[string]interface{})["id"].(float64)) == chall.ID {
			challengeBody = v
			break
		}
	}
	tags := challengeBody.(map[string]interface{})["tags"].([]interface{})
	if len(tags) != 0 {
		t.Fatalf("Expected no tags, but got: %v", tags)
	}
}
